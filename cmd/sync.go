package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/Supme/uiscom"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Version = "0.1.0"
var verbose bool

func main() {
	var (
		uiscomToken string

		dbHost     string
		dbPort     int
		dbName     string
		dbUser     string
		dbPassword string

		fromStr string

		mediaFolder string
	)

	flag.StringVar(&uiscomToken, "t", "", "Uiscom user token")

	flag.StringVar(&dbHost, "s", "", "Database host")
	flag.IntVar(&dbPort, "p", 5432, "Database port")
	flag.StringVar(&dbName, "n", "", "Database name")
	flag.StringVar(&dbUser, "u", "", "Database user")
	flag.StringVar(&dbPassword, "w", "", "Database user password")

	interval := flag.Duration("i", time.Hour, "Interval days ago (eg: 4h, 60m, 45s, 12h15m30s)")
	flag.StringVar(&fromStr, "f", "", "From datetime, format \""+uiscom.DateFormat+"\", default (time now - interval)")

	flag.StringVar(&mediaFolder, "m", "", "Folder for sync media records (blanc not syncing)")

	flag.BoolVar(&verbose, "V", false, "Verbose output")

	version := flag.Bool("v", false, "Prints version")

	flag.Parse()

	if *version {
		fmt.Printf("Uiscom syncronize to Postgresql version: v%s\r\n\r\n", Version)
		os.Exit(0)
	}

	var till, from time.Time
	if fromStr == "" {
		till = time.Now()
		from = till.Add(-*interval)
	} else {
		var err error
		from, err = uiscom.StringToTime(fromStr)
		if err != nil {
			fmt.Printf("Wrong from datetime format: %s\r\n", err)
			os.Exit(1)
		}
		till = from.Add(*interval)
	}

	if verbose {
		log.Printf("Start syncronize between \"%s\" and \"%s\"\r\n", uiscom.TimeToString(from), uiscom.TimeToString(till))
	}

	dburl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	dbpool, err := pgxpool.New(context.Background(), dburl)
	//db, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	if verbose {
		log.Print("database connected")
	}

	client := uiscom.NewWithToken(uiscom.TargetUiscom, uiscomToken)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if verbose {
			log.Print("start calls syncing")
		}
		err := syncCalls(dbpool, client, from, till, mediaFolder)
		if err != nil {
			log.Printf("error calls syncing %s", err)
		}
		if verbose {
			log.Print("finish calls syncing")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if verbose {
			log.Print("start call_legs syncing")
		}
		err := syncCallLegs(dbpool, client, from, till)
		if err != nil {
			log.Printf("error call_legs syncing %s", err)
		}
		if verbose {
			log.Print("finish call_legs syncing")
		}
	}()

	wg.Wait()
	if verbose {
		log.Print("all finish")
	}
}

func syncCalls(dbpool *pgxpool.Pool, client *uiscom.Client, from, till time.Time, mediaFolder string) error {
	fields := append(
		uiscom.GetCallsReportResponseParametersFields,
		append(
			uiscom.GetCallsReportResponseFirstAnsweredEmployeeFields,
			append(
				uiscom.GetCallsReportResponseLastAnsweredEmployeeFields,
				append(
					uiscom.GetCallsReportResponseScenarioFields,
					append(
						uiscom.GetCallsReportResponseFirstTalkedEmployeeFields,
						[]uiscom.Field{
							"communication_id",
							"call_records",
							//"wav_call_records", // всегда пустой массив []
							//"full_record_file_link", // ссылка на склеенную запись
							"voice_mail_records",
							"virtual_phone_number",
							"source",
							"is_lost",
							"contact_phone_number",
						}...,
					)...,
				)...,
			)...,
		)...,
	)

	resp, err := client.GetCalls(context.Background(), -1, from, till, 10000, 0, nil, fields...)
	if err != nil {
		return err
	}

	if r, ok := resp.(map[string]any)["data"]; ok {
		for _, v := range r.([]any) {
			val, err := validate(fields, v)
			if err != nil {
				return err
			}

			errDownload := make(chan error)
			go func() {
				if mediaFolder != "" {
					var mediaURLs []string
					var communicationFolder string

					if len(val["call_records"].([]any)) != 0 {
						for _, r := range val["call_records"].([]any) {
							u, err := url.JoinPath(uiscom.UiscomTalkMediaURL,
								strconv.FormatInt(val["communication_id"].(int64), 10),
								r.(string), "/")
							if err != nil {
								errDownload <- err
								return
							}
							mediaURLs = append(mediaURLs, u)
						}
						communicationFolder = strconv.FormatInt(val["communication_id"].(int64), 10)
						if val["direction"].(string) == "out" {
							communicationFolder = "out_" + communicationFolder
						}
					} else if len(val["voice_mail_records"].([]any)) != 0 {
						for _, r := range val["voice_mail_records"].([]any) {
							u, err := url.JoinPath(uiscom.UiscomVoiceMailMediaURL,
								strconv.FormatInt(val["communication_id"].(int64), 10),
								r.(string), "/")
							if err != nil {
								errDownload <- err
								return
							}
							mediaURLs = append(mediaURLs, u)
						}
						communicationFolder = "vm_" + strconv.FormatInt(val["communication_id"].(int64), 10)
					}

					mFolder := filepath.Join(mediaFolder,
						fmt.Sprintf("%4d", val["start_time"].(time.Time).Year()),
						fmt.Sprintf("%02d", val["start_time"].(time.Time).Month()),
						fmt.Sprintf("%02d", val["start_time"].(time.Time).Day()),
						communicationFolder,
					) + "/"

					err = RecordsDownload(mFolder, mediaURLs...)
					if err != nil {
						errDownload <- err
						return
					}
				}
				errDownload <- nil
				return
			}()

			_, err = dbpool.Exec(context.Background(),
				`INSERT INTO calls
				(id, communication_id, start_time, finish_time, finish_reason, direction, is_lost, virtual_phone_number, contact_phone_number, first_answered_employee_id, first_answered_employee_full_name, first_talked_employee_id, first_talked_employee_full_name, last_answered_employee_id, last_answered_employee_full_name, scenario_id, scenario_name, source)
				VALUES
				($1, $2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
				ON CONFLICT DO NOTHING`,
				val["id"],
				val["communication_id"],
				val["start_time"],
				val["finish_time"],
				val["finish_reason"],
				val["direction"],
				val["is_lost"],
				val["virtual_phone_number"],
				val["contact_phone_number"],
				val["first_answered_employee_id"],
				val["first_answered_employee_full_name"],
				val["first_talked_employee_id"],
				val["first_talked_employee_full_name"],
				val["last_answered_employee_id"],
				val["last_answered_employee_full_name"],
				val["scenario_id"],
				val["scenario_name"],
				val["source"],
			)
			if err != nil {
				return err
			}

			err = <-errDownload
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("response not have field 'data'")
	}

	return nil
}

func syncCallLegs(dbpool *pgxpool.Pool, client *uiscom.Client, from, till time.Time) error {
	fields := uiscom.GetCallLegsReportResponseParametersFields

	resp, err := client.GetCallLegs(context.Background(), -1, from, till, 10000, 0, nil, fields...)
	if err != nil {
		return err
	}

	if r, ok := resp.(map[string]any)["data"]; ok {
		for _, v := range r.([]any) {
			val, err := validate(fields, v)
			if err != nil {
				return err
			}

			_, err = dbpool.Exec(context.Background(),
				`INSERT INTO call_legs (
				    id,
    				call_session_id,
    				start_time,
    				connect_time,
    				duration,
    				total_duration,
    				finish_reason,
    				finish_reason_description,
    				virtual_phone_number,
    				calling_phone_number,
    				called_phone_number,
    				direction,
    				is_transfered,
    				is_operator,
    				is_coach,
    				is_failed,
    				is_talked,
    				employee_id,
    				employee_full_name,
    				employee_phone_number,
    				scenario_id,
    				scenario_name,
    				release_cause_code,
    				release_cause_description,
    				contact_id,
    				contact_full_name,
    				contact_phone_number,
    				action_id,
    				action_name,
    				group_id,
    				group_name)
				VALUES
					($1, $2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31)
				ON CONFLICT DO NOTHING`,
				val["id"],
				val["call_session_id"],
				val["start_time"],
				val["connect_time"],
				val["duration"].(time.Duration),
				durationToInterval(val["total_duration"].(time.Duration)),
				val["finish_reason"],
				val["finish_reason_description"],
				val["virtual_phone_number"],
				val["calling_phone_number"],
				val["called_phone_number"],
				val["direction"],
				val["is_transfered"],
				val["is_operator"],
				val["is_coach"],
				val["is_failed"],
				val["is_talked"],
				val["employee_id"],
				val["employee_full_name"],
				val["employee_phone_number"],
				val["scenario_id"],
				val["scenario_name"],
				val["release_cause_code"],
				val["release_cause_description"],
				val["contact_id"],
				val["contact_full_name"],
				val["contact_phone_number"],
				val["action_id"],
				val["action_name"],
				val["group_id"],
				val["group_name"],
			)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("response not have field 'data'")
	}

	return nil
}

func validate(fields []uiscom.Field, data any) (map[string]any, error) {
	val := make(map[string]any)
	absent := []string{}
	for i := range fields {
		if _, ok := data.(map[string]any)[fields[i].String()]; !ok {
			absent = append(absent, fields[i].String())
		} else {
			if data.(map[string]any)[fields[i].String()] == nil {
				val[fields[i].String()] = nil
			} else {
				switch reflect.TypeOf(data.(map[string]any)[fields[i].String()]).Name() {
				case "Number":
					switch {
					case strings.HasSuffix(fields[i].String(), "duration"):
						var err error
						val[fields[i].String()], err = time.ParseDuration(data.(map[string]any)[fields[i].String()].(json.Number).String() + "s")
						if err != nil {
							return val, err
						}
					default:
						val[fields[i].String()], _ = data.(map[string]any)[fields[i].String()].(json.Number).Int64()
					}

				case "string":
					switch {
					case strings.HasSuffix(fields[i].String(), "time"):
						var err error
						val[fields[i].String()], err = uiscom.StringToTime(data.(map[string]any)[fields[i].String()].(string))
						if err != nil {
							return val, err
						}
					default:
						val[fields[i].String()] = data.(map[string]any)[fields[i].String()].(string)
					}

				default:
					val[fields[i].String()] = data.(map[string]any)[fields[i].String()]
				}
			}
		}
	}
	if len(absent) != 0 {
		return val, fmt.Errorf("absent %v fields", absent)
	}

	return val, nil
}

func durationToInterval(duration time.Duration) pgtype.Interval {
	interval := pgtype.Interval{}
	_ = interval.Set(duration)
	return interval
}

func RecordsDownload(mediaFolder string, urls ...string) error {
	for i := range urls {
		err := download(mediaFolder, urls[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func download(folder, url string) error {
	if folder == "" {
		folder = "."
	}
	if _, err := os.Stat(folder); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(folder, 0755); err != nil {
				return errors.New("error creating destination directory: " + err.Error())
			}
		} else {
			return errors.New("error checking destination directory: " + err.Error())
		}
	}

	client := http.Client{}
	request, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	filename := resp.Request.URL.Path
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if _, params, err := mime.ParseMediaType(cd); err == nil {
			if val, ok := params["filename"]; ok {
				filename = val
			}
		}
	}
	filename = filepath.Base(path.Clean("/" + filename))
	if filename == "" || filename == "." || filename == "/" {
		return errors.New("filename couldn't be determined")
	}
	filename = filepath.Join(folder, filename)

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			if verbose {
				fmt.Println("download and create", filename)
			}
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, resp.Body)
			if err != nil {
				return err
			}
		} else {
			return errors.New("error checking destination file: " + err.Error())
		}
	} else {
		if verbose {
			fmt.Println("file", filename, "already exist")
		}
	}
	return nil
}
