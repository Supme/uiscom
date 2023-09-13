package uiscom

import (
	"fmt"
	"strconv"
)

type FilterCondition string

func (c FilterCondition) String() string {
	return string(c)
}

const (
	FilterConditionOr   = FilterCondition("or")
	FilterConditionAnd  = FilterCondition("and")
	FilterConditionNone = FilterCondition("")
)

type Filter struct {
	Field    string
	Operator string
	Value    any

	filters   []Filter
	condition FilterCondition
}

func (f Filter) JsonPart() string {
	if len(f.filters) == 0 {
		return fmt.Sprintf("{\"field\":%s,\"operator\":%s,\"value\":%s}", strconv.Quote(f.Field), "\""+f.Operator+"\"", normalizeJsonValue(f.Value))
	} else if len(f.filters) == 1 {
		return f.filters[0].JsonPart()
	}

	s := "{\"filters\":["
	for i := range f.filters {
		s += f.filters[i].JsonPart()
		if i != len(f.filters)-1 {
			s += ","
		}
	}
	s += "],"
	if f.condition != "" {
		s += "\"condition\": \"" + string(f.condition) + "\"}"
	}
	return s
}

func normalizeJsonValue(v any) string {
	switch v.(type) {
	case string:
		return strconv.Quote(v.(string))
	case int:
		return strconv.Itoa(v.(int))
	case bool:
		if v.(bool) {
			return "true"
		} else {
			return "false"
		}
	case nil:
		return "null"
	default:
		return "unsupported type"
	}
}

func GetFilterSingle(filter Filter) *Filter {
	return filterSet(FilterConditionNone, filter)
}

func GetFilterOr(filter ...Filter) *Filter {
	return filterSet(FilterConditionOr, filter...)
}

func GetFilterAnd(filter ...Filter) *Filter {
	return filterSet(FilterConditionAnd, filter...)
}

func filterSet(condition FilterCondition, filter ...Filter) *Filter {
	var f Filter
	f.filters = append(f.filters, filter...)
	f.condition = condition
	return &f
}

type SortOrder string

const (
	SortOrderAsc  = SortOrder("asc")
	SortOrderDesc = SortOrder("desc")
)

func (o SortOrder) String() string {
	return string(o)
}

type Sort struct {
	Field string
	Order SortOrder
}

type Field string

func (f Field) String() string {
	return string(f)
}

// Параметры ответа
var GetCallsReportResponseParametersFields = []Field{
	"id",
	"start_time",
	"finish_time",
	"finish_reason",
	"direction",
	"cpn_region_id",
	"cpn_region_name",
}

// Операции сценариев
var GetCallsReportResponseParametersScriptOperationsFields = []Field{
	"scenario_operations",
	"id",
	"name",
	"source",
	"is_lost",
	"communication_number",
	"communication_page_url",
	"contact_phone_number",
	"communication_id", // всегда == id ?
	"communication_type",
	"wait_duration",
	"total_wait_duration",
	"lost_call_processing_duration",
	"talk_duration",
	"clean_talk_duration",
	"total_duration",
	"postprocess_duration",
	"call_records", // https://app.uiscom.ru/system/media/talk/{communication_id}/{[call_records]}/
	"wav_call_records",
	"full_record_file_link", // ссылка на склееную запись https://media.uiscom.ru/{communication_id}/{call_records[0]}
	"voice_mail_records",
	"virtual_phone_number",
	"ua_client_id",
	"ym_client_id",
	"sale_date",
	"sale_cost",
	"is_transfer",
	"search_query",
	"search_engine",
	"referrer_domain",
	"referrer",
	"entrance_page",
	"gclid",
	"yclid",
	"ymclid",
	"ef_id",
	"channel",
}

// Проставленные теги
var GetCallsReportResponseAttachedTagsFields = []Field{
	"tags",
	"tag_name",
	"tag_id",
	"tag_change_time",
	"tag_type",
	"tag_user_id",
	"tag_user_login",
	"tag_employee_id",
	"tag_employee_full_name",
}

// Сотрудники участвовавшие в звонке
var GetCallsReportResponseEmployeesParticipatedInCallFields = []Field{
	"employees",
	"employee_id",
	"employee_full_name",
	"is_answered",
	"is_talked",
}

// Последний ответивший сотрудник
var GetCallsReportResponseLastAnsweredEmployeeFields = []Field{
	"last_answered_employee_id",
	"last_answered_employee_full_name",
	"last_answered_employee_rating",
}

// Первый ответивший сотрудник
var GetCallsReportResponseFirstAnsweredEmployeeFields = []Field{
	"first_answered_employee_id",
	"first_answered_employee_full_name",
}

// Последний разговаривавший сотрудник
var GetCallsReportResponseFirstTalkedEmployeeFields = []Field{
	"first_talked_employee_id",
	"first_talked_employee_full_name",
}

// Сценарий
var GetCallsReportResponseScenarioFields = []Field{
	"scenario_name",
	"scenario_id",
}

// Сайт
var GetCallsReportResponseSiteFields = []Field{
	"site_domain_name",
	"site_id",
}

// Рекламная кампания
var GetCallsReportResponseCampaignFields = []Field{
	"campaign_name",
	"campaign_id",
	"visit_other_campaign",
}

// Информация о посетителе
var GetCallsReportResponseVisitorFields = []Field{
	"visitor_id",
	"person_id",
	"visitor_type",
	"visitor_session_id",
	"visits_count",
	"visitor_first_campaign_id",
	"visitor_first_campaign_name",
	"visitor_city",
	"visitor_region",
	"visitor_country",
	"visitor_device",
}

// Свойства посетителя
var GetCallsReportResponseVisitorPropertiesFields = []Field{
	"visitor_custom_properties",
	"property_name",
	"property_value",
}

// Сегменты
var GetCallsReportResponseSegmentsFields = []Field{
	"segments",
	"segment_id",
	"segment_name",
}

// Call API
var GetCallsReportResponseCallApiFields = []Field{
	"call_api_request_id",
	"call_api_external_id",
}

// Контакт из адресной книги
var GetCallsReportResponseContactFields = []Field{
	"contact_id",
	"contact_full_name",
}

// UTM-метки
var GetCallsReportResponseUtmFields = []Field{
	"utm_source",
	"utm_medium",
	"utm_term",
	"utm_content",
	"utm_campaign",
}

// OS-метки
var GetCallsReportResponseOpenstatFields = []Field{
	"openstat_ad",
	"openstat_campaign",
	"openstat_service",
	"openstat_source",
}

// Атрибуты обращения
var GetCallsReportResponseAtributesFields = []Field{
	"attributes",
}

// Расширенные UTM-метки
var GetCallsReportResponseEqUtmFields = []Field{
	"eq_utm_source",
	"eq_utm_medium",
	"eq_utm_term",
	"eq_utm_content",
	"eq_utm_campaign",
	"eq_utm_referrer",
	"eq_utm_expid",
}

var GetCallLegsReportResponseParametersFields = []Field{
	"id",
	"call_session_id",
	"call_records",
	"wav_call_records",
	"start_time",
	"connect_time",
	"duration",
	"total_duration",
	"finish_reason",
	"finish_reason_description",
	"virtual_phone_number",
	"calling_phone_number",
	"called_phone_number",
	"direction",
	"is_transfered",
	"is_operator",
	"employee_id",
	"employee_full_name",
	"employee_phone_number",
	"employee_rating",
	"scenario_id",
	"scenario_name",
	"is_coach",
	"release_cause_code",
	"release_cause_description",
	"is_failed",
	"is_talked",
	"contact_id",
	"contact_full_name",
	"contact_phone_number",
	"action_id",
	"action_name",
	"group_id",
	"group_name",
}
