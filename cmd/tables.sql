-- Table: public.calls

-- DROP TABLE IF EXISTS public.calls;

CREATE TABLE IF NOT EXISTS public.calls
(
    id bigint NOT NULL,
    communication_id bigint,
    start_time timestamp without time zone,
    finish_time timestamp without time zone,
    finish_reason character varying(100) COLLATE pg_catalog."default",
    direction character varying(3) COLLATE pg_catalog."default",
    is_lost boolean,
    virtual_phone_number character varying(20) COLLATE pg_catalog."default",
    contact_phone_number character varying(20) COLLATE pg_catalog."default",
    first_answered_employee_id bigint,
    first_answered_employee_full_name character varying(250) COLLATE pg_catalog."default",
    first_talked_employee_id bigint,
    first_talked_employee_full_name character varying(250) COLLATE pg_catalog."default",
    last_answered_employee_id bigint,
    last_answered_employee_full_name character varying(250) COLLATE pg_catalog."default",
    scenario_id bigint,
    scenario_name character varying(250) COLLATE pg_catalog."default",
    source character varying(100) COLLATE pg_catalog."default",
    CONSTRAINT calls_pkey PRIMARY KEY (id)
    )

-- Table: public.call_legs

-- DROP TABLE IF EXISTS public.call_legs;

CREATE TABLE IF NOT EXISTS public.call_legs
(
    id bigint NOT NULL,
    call_session_id bigint,
    start_time timestamp without time zone,
    connect_time timestamp without time zone,
    duration interval,
    total_duration interval,
    finish_reason character varying(100) COLLATE pg_catalog."default",
    finish_reason_description character varying(100) COLLATE pg_catalog."default",
    virtual_phone_number character varying(20) COLLATE pg_catalog."default",
    calling_phone_number character varying(20) COLLATE pg_catalog."default",
    called_phone_number character varying(20) COLLATE pg_catalog."default",
    direction character varying(3) COLLATE pg_catalog."default",
    is_transfered boolean,
    is_operator boolean,
    is_coach boolean,
    is_failed boolean,
    is_talked boolean,
    employee_id bigint,
    employee_full_name character varying(250) COLLATE pg_catalog."default",
    employee_phone_number character varying(20) COLLATE pg_catalog."default",
    scenario_id bigint,
    scenario_name character varying(250) COLLATE pg_catalog."default",
    release_cause_code bigint,
    release_cause_description character varying(250) COLLATE pg_catalog."default",
    contact_id bigint,
    contact_full_name character varying(250) COLLATE pg_catalog."default",
    contact_phone_number character varying(20) COLLATE pg_catalog."default",
    action_id bigint,
    action_name character varying(250) COLLATE pg_catalog."default",
    group_id bigint,
    group_name character varying(250) COLLATE pg_catalog."default",
    CONSTRAINT call_legs_pkey PRIMARY KEY (id)
    )



ALTER DEFAULT PRIVILEGES FOR ROLE uiscom
GRANT SELECT ON TABLES TO uiscom_reader;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres
GRANT ALL ON TABLES TO uiscom;

ALTER DEFAULT PRIVILEGES FOR ROLE "a.agafonov"
GRANT SELECT ON TABLES TO uiscom_reader;

ALTER DEFAULT PRIVILEGES FOR ROLE nkozkin
GRANT SELECT ON TABLES TO uiscom_reader;

ALTER DEFAULT PRIVILEGES FOR ROLE "a.agafonov"
GRANT ALL ON TABLES TO uiscom;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres
GRANT SELECT ON TABLES TO uiscom_reader;

ALTER DEFAULT PRIVILEGES FOR ROLE nkozkin
GRANT ALL ON TABLES TO uiscom;



