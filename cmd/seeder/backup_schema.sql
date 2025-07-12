--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5
-- Dumped by pg_dump version 17.5

-- Started on 2025-07-12 14:44:35

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 249 (class 1255 OID 16609)
-- Name: trigger_set_timestamp(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.trigger_set_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;$$;


ALTER FUNCTION public.trigger_set_timestamp() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 223 (class 1259 OID 16420)
-- Name: area; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.area (
    id smallint NOT NULL,
    department_id smallint NOT NULL,
    name text NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.area OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 16419)
-- Name: area_department_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.area_department_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.area_department_id_seq OWNER TO postgres;

--
-- TOC entry 3621 (class 0 OID 0)
-- Dependencies: 222
-- Name: area_department_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.area_department_id_seq OWNED BY public.area.department_id;


--
-- TOC entry 221 (class 1259 OID 16418)
-- Name: area_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.area_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.area_id_seq OWNER TO postgres;

--
-- TOC entry 3622 (class 0 OID 0)
-- Dependencies: 221
-- Name: area_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.area_id_seq OWNED BY public.area.id;


--
-- TOC entry 220 (class 1259 OID 16404)
-- Name: department; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.department (
    id smallint NOT NULL,
    name text NOT NULL,
    receive_job boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.department OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 16403)
-- Name: department_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.department_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.department_id_seq OWNER TO postgres;

--
-- TOC entry 3623 (class 0 OID 0)
-- Dependencies: 219
-- Name: department_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.department_id_seq OWNED BY public.department.id;


--
-- TOC entry 226 (class 1259 OID 16441)
-- Name: employee; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.employee (
    npk text NOT NULL,
    department_id smallint,
    area_id smallint,
    name text NOT NULL,
    "position" text NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.employee OWNER TO postgres;

--
-- TOC entry 225 (class 1259 OID 16440)
-- Name: employee_area_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.employee_area_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.employee_area_id_seq OWNER TO postgres;

--
-- TOC entry 3624 (class 0 OID 0)
-- Dependencies: 225
-- Name: employee_area_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.employee_area_id_seq OWNED BY public.employee.area_id;


--
-- TOC entry 224 (class 1259 OID 16439)
-- Name: employee_department_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.employee_department_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.employee_department_id_seq OWNER TO postgres;

--
-- TOC entry 3625 (class 0 OID 0)
-- Dependencies: 224
-- Name: employee_department_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.employee_department_id_seq OWNED BY public.employee.department_id;


--
-- TOC entry 243 (class 1259 OID 16562)
-- Name: job; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.job (
    id bigint NOT NULL,
    ticket_id bigint NOT NULL,
    pic_job text,
    job_priority bigint NOT NULL,
    report_file text[],
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.job OWNER TO postgres;

--
-- TOC entry 241 (class 1259 OID 16560)
-- Name: job_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.job_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.job_id_seq OWNER TO postgres;

--
-- TOC entry 3626 (class 0 OID 0)
-- Dependencies: 241
-- Name: job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.job_id_seq OWNED BY public.job.id;


--
-- TOC entry 242 (class 1259 OID 16561)
-- Name: job_ticket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.job_ticket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.job_ticket_id_seq OWNER TO postgres;

--
-- TOC entry 3627 (class 0 OID 0)
-- Dependencies: 242
-- Name: job_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.job_ticket_id_seq OWNED BY public.job.ticket_id;


--
-- TOC entry 228 (class 1259 OID 16464)
-- Name: physical_location; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.physical_location (
    id smallint NOT NULL,
    name text NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.physical_location OWNER TO postgres;

--
-- TOC entry 227 (class 1259 OID 16463)
-- Name: physical_location_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.physical_location_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.physical_location_id_seq OWNER TO postgres;

--
-- TOC entry 3628 (class 0 OID 0)
-- Dependencies: 227
-- Name: physical_location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.physical_location_id_seq OWNED BY public.physical_location.id;


--
-- TOC entry 246 (class 1259 OID 16587)
-- Name: rejected_ticket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rejected_ticket (
    id bigint NOT NULL,
    ticket_id bigint NOT NULL,
    rejector text NOT NULL,
    feedback text NOT NULL,
    already_seen boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.rejected_ticket OWNER TO postgres;

--
-- TOC entry 244 (class 1259 OID 16585)
-- Name: rejected_ticket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.rejected_ticket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.rejected_ticket_id_seq OWNER TO postgres;

--
-- TOC entry 3629 (class 0 OID 0)
-- Dependencies: 244
-- Name: rejected_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.rejected_ticket_id_seq OWNED BY public.rejected_ticket.id;


--
-- TOC entry 245 (class 1259 OID 16586)
-- Name: rejected_ticket_ticket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.rejected_ticket_ticket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.rejected_ticket_ticket_id_seq OWNER TO postgres;

--
-- TOC entry 3630 (class 0 OID 0)
-- Dependencies: 245
-- Name: rejected_ticket_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.rejected_ticket_ticket_id_seq OWNED BY public.rejected_ticket.ticket_id;


--
-- TOC entry 231 (class 1259 OID 16479)
-- Name: specified_location; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.specified_location (
    id smallint NOT NULL,
    physical_location_id smallint NOT NULL,
    name text NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.specified_location OWNER TO postgres;

--
-- TOC entry 229 (class 1259 OID 16477)
-- Name: specified_location_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.specified_location_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.specified_location_id_seq OWNER TO postgres;

--
-- TOC entry 3631 (class 0 OID 0)
-- Dependencies: 229
-- Name: specified_location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.specified_location_id_seq OWNED BY public.specified_location.id;


--
-- TOC entry 230 (class 1259 OID 16478)
-- Name: specified_location_physical_location_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.specified_location_physical_location_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.specified_location_physical_location_id_seq OWNER TO postgres;

--
-- TOC entry 3632 (class 0 OID 0)
-- Dependencies: 230
-- Name: specified_location_physical_location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.specified_location_physical_location_id_seq OWNED BY public.specified_location.physical_location_id;


--
-- TOC entry 218 (class 1259 OID 16390)
-- Name: status_ticket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.status_ticket (
    id smallint NOT NULL,
    name text NOT NULL,
    sequence smallint NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.status_ticket OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16389)
-- Name: status_ticket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.status_ticket_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.status_ticket_id_seq OWNER TO postgres;

--
-- TOC entry 3633 (class 0 OID 0)
-- Dependencies: 217
-- Name: status_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.status_ticket_id_seq OWNED BY public.status_ticket.id;


--
-- TOC entry 236 (class 1259 OID 16502)
-- Name: ticket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ticket (
    id bigint NOT NULL,
    requestor text NOT NULL,
    department_target_id smallint NOT NULL,
    physical_location_id smallint,
    specified_location_id smallint,
    description text NOT NULL,
    ticket_priority bigint NOT NULL,
    support_file text[],
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.ticket OWNER TO postgres;

--
-- TOC entry 233 (class 1259 OID 16499)
-- Name: ticket_department_target_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ticket_department_target_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ticket_department_target_id_seq OWNER TO postgres;

--
-- TOC entry 3634 (class 0 OID 0)
-- Dependencies: 233
-- Name: ticket_department_target_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ticket_department_target_id_seq OWNED BY public.ticket.department_target_id;


--
-- TOC entry 232 (class 1259 OID 16498)
-- Name: ticket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ticket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ticket_id_seq OWNER TO postgres;

--
-- TOC entry 3635 (class 0 OID 0)
-- Dependencies: 232
-- Name: ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ticket_id_seq OWNED BY public.ticket.id;


--
-- TOC entry 234 (class 1259 OID 16500)
-- Name: ticket_physical_location_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ticket_physical_location_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ticket_physical_location_id_seq OWNER TO postgres;

--
-- TOC entry 3636 (class 0 OID 0)
-- Dependencies: 234
-- Name: ticket_physical_location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ticket_physical_location_id_seq OWNED BY public.ticket.physical_location_id;


--
-- TOC entry 235 (class 1259 OID 16501)
-- Name: ticket_specified_location_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ticket_specified_location_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ticket_specified_location_id_seq OWNER TO postgres;

--
-- TOC entry 3637 (class 0 OID 0)
-- Dependencies: 235
-- Name: ticket_specified_location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ticket_specified_location_id_seq OWNED BY public.ticket.specified_location_id;


--
-- TOC entry 240 (class 1259 OID 16538)
-- Name: track_status_ticket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.track_status_ticket (
    id bigint NOT NULL,
    ticket_id bigint NOT NULL,
    status_ticket_id smallint NOT NULL,
    start_date timestamp with time zone DEFAULT now() NOT NULL,
    finish_date timestamp with time zone
);


ALTER TABLE public.track_status_ticket OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 16535)
-- Name: track_status_ticket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.track_status_ticket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.track_status_ticket_id_seq OWNER TO postgres;

--
-- TOC entry 3638 (class 0 OID 0)
-- Dependencies: 237
-- Name: track_status_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.track_status_ticket_id_seq OWNED BY public.track_status_ticket.id;


--
-- TOC entry 239 (class 1259 OID 16537)
-- Name: track_status_ticket_status_ticket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.track_status_ticket_status_ticket_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.track_status_ticket_status_ticket_id_seq OWNER TO postgres;

--
-- TOC entry 3639 (class 0 OID 0)
-- Dependencies: 239
-- Name: track_status_ticket_status_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.track_status_ticket_status_ticket_id_seq OWNED BY public.track_status_ticket.status_ticket_id;


--
-- TOC entry 238 (class 1259 OID 16536)
-- Name: track_status_ticket_ticket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.track_status_ticket_ticket_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.track_status_ticket_ticket_id_seq OWNER TO postgres;

--
-- TOC entry 3640 (class 0 OID 0)
-- Dependencies: 238
-- Name: track_status_ticket_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.track_status_ticket_ticket_id_seq OWNED BY public.track_status_ticket.ticket_id;


--
-- TOC entry 248 (class 1259 OID 16624)
-- Name: view_calculation_ticket_each_status; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.view_calculation_ticket_each_status AS
 WITH ticket_cohorts AS (
         SELECT tst.ticket_id,
            min(tst.start_date) AS cohort_date
           FROM (public.track_status_ticket tst
             JOIN public.status_ticket st ON ((tst.status_ticket_id = st.id)))
          WHERE (st.sequence >= 0)
          GROUP BY tst.ticket_id
        ), ranked_ticket_statuses AS (
         SELECT tst.ticket_id,
            tst.status_ticket_id,
            row_number() OVER (PARTITION BY tst.ticket_id ORDER BY tst.start_date DESC, tst.id DESC) AS rank_num
           FROM public.track_status_ticket tst
          WHERE (tst.ticket_id IN ( SELECT ticket_cohorts.ticket_id
                   FROM ticket_cohorts))
        ), current_ticket_status AS (
         SELECT tc.ticket_id,
            tc.cohort_date,
            st.name AS current_status_name
           FROM ((ticket_cohorts tc
             JOIN ranked_ticket_statuses rts ON (((tc.ticket_id = rts.ticket_id) AND (rts.rank_num = 1))))
             JOIN public.status_ticket st ON ((rts.status_ticket_id = st.id)))
        ), distinct_periods AS (
         SELECT DISTINCT (EXTRACT(year FROM ticket_cohorts.cohort_date))::integer AS year,
            (EXTRACT(month FROM ticket_cohorts.cohort_date))::integer AS month
           FROM ticket_cohorts
        ), relevant_statuses AS (
         SELECT status_ticket.name
           FROM public.status_ticket
          WHERE ((status_ticket.sequence >= 0) OR (status_ticket.sequence = '-100'::integer))
        )
 SELECT dp.year,
    dp.month,
    rs.name AS status,
    COALESCE(count(cts.ticket_id), (0)::bigint) AS total
   FROM ((distinct_periods dp
     CROSS JOIN relevant_statuses rs)
     LEFT JOIN current_ticket_status cts ON ((((dp.year)::numeric = EXTRACT(year FROM cts.cohort_date)) AND ((dp.month)::numeric = EXTRACT(month FROM cts.cohort_date)) AND (rs.name = cts.current_status_name))))
  GROUP BY dp.year, dp.month, rs.name
  ORDER BY dp.year, dp.month, ( SELECT status_ticket.sequence
           FROM public.status_ticket
          WHERE (status_ticket.name = rs.name));


ALTER VIEW public.view_calculation_ticket_each_status OWNER TO postgres;

--
-- TOC entry 247 (class 1259 OID 16619)
-- Name: view_ticket_list; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.view_ticket_list AS
 SELECT t.id AS ticket_id,
    t.description,
    t.ticket_priority,
    loc.name AS location_name,
    spec_loc.name AS specified_location_name,
    ((now())::date - (t.created_at)::date) AS ticket_age_days,
    req_emp.name AS requestor_name,
    req_dept.name AS requestor_department,
    ( SELECT st.name
           FROM (public.track_status_ticket tst
             JOIN public.status_ticket st ON ((tst.status_ticket_id = st.id)))
          WHERE ((tst.ticket_id = t.id) AND (tst.finish_date IS NULL))
         LIMIT 1) AS current_status,
    pic_emp.name AS pic_name,
    pic_area.name AS pic_area_name
   FROM (((((((public.ticket t
     JOIN public.employee req_emp ON ((t.requestor = req_emp.npk)))
     LEFT JOIN public.department req_dept ON ((req_emp.department_id = req_dept.id)))
     LEFT JOIN public.physical_location loc ON ((t.physical_location_id = loc.id)))
     LEFT JOIN public.specified_location spec_loc ON ((t.specified_location_id = spec_loc.id)))
     LEFT JOIN public.job j ON ((t.id = j.ticket_id)))
     LEFT JOIN public.employee pic_emp ON ((j.pic_job = pic_emp.npk)))
     LEFT JOIN public.area pic_area ON ((pic_emp.area_id = pic_area.id)));


ALTER VIEW public.view_ticket_list OWNER TO postgres;

--
-- TOC entry 3344 (class 2604 OID 16423)
-- Name: area id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.area ALTER COLUMN id SET DEFAULT nextval('public.area_id_seq'::regclass);


--
-- TOC entry 3345 (class 2604 OID 16424)
-- Name: area department_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.area ALTER COLUMN department_id SET DEFAULT nextval('public.area_department_id_seq'::regclass);


--
-- TOC entry 3339 (class 2604 OID 16407)
-- Name: department id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.department ALTER COLUMN id SET DEFAULT nextval('public.department_id_seq'::regclass);


--
-- TOC entry 3349 (class 2604 OID 16444)
-- Name: employee department_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee ALTER COLUMN department_id SET DEFAULT nextval('public.employee_department_id_seq'::regclass);


--
-- TOC entry 3350 (class 2604 OID 16445)
-- Name: employee area_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee ALTER COLUMN area_id SET DEFAULT nextval('public.employee_area_id_seq'::regclass);


--
-- TOC entry 3373 (class 2604 OID 16565)
-- Name: job id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job ALTER COLUMN id SET DEFAULT nextval('public.job_id_seq'::regclass);


--
-- TOC entry 3374 (class 2604 OID 16566)
-- Name: job ticket_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job ALTER COLUMN ticket_id SET DEFAULT nextval('public.job_ticket_id_seq'::regclass);


--
-- TOC entry 3354 (class 2604 OID 16467)
-- Name: physical_location id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.physical_location ALTER COLUMN id SET DEFAULT nextval('public.physical_location_id_seq'::regclass);


--
-- TOC entry 3377 (class 2604 OID 16590)
-- Name: rejected_ticket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rejected_ticket ALTER COLUMN id SET DEFAULT nextval('public.rejected_ticket_id_seq'::regclass);


--
-- TOC entry 3378 (class 2604 OID 16591)
-- Name: rejected_ticket ticket_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rejected_ticket ALTER COLUMN ticket_id SET DEFAULT nextval('public.rejected_ticket_ticket_id_seq'::regclass);


--
-- TOC entry 3358 (class 2604 OID 16482)
-- Name: specified_location id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.specified_location ALTER COLUMN id SET DEFAULT nextval('public.specified_location_id_seq'::regclass);


--
-- TOC entry 3359 (class 2604 OID 16483)
-- Name: specified_location physical_location_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.specified_location ALTER COLUMN physical_location_id SET DEFAULT nextval('public.specified_location_physical_location_id_seq'::regclass);


--
-- TOC entry 3335 (class 2604 OID 16393)
-- Name: status_ticket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.status_ticket ALTER COLUMN id SET DEFAULT nextval('public.status_ticket_id_seq'::regclass);


--
-- TOC entry 3363 (class 2604 OID 16505)
-- Name: ticket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket ALTER COLUMN id SET DEFAULT nextval('public.ticket_id_seq'::regclass);


--
-- TOC entry 3364 (class 2604 OID 16506)
-- Name: ticket department_target_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket ALTER COLUMN department_target_id SET DEFAULT nextval('public.ticket_department_target_id_seq'::regclass);


--
-- TOC entry 3365 (class 2604 OID 16507)
-- Name: ticket physical_location_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket ALTER COLUMN physical_location_id SET DEFAULT nextval('public.ticket_physical_location_id_seq'::regclass);


--
-- TOC entry 3366 (class 2604 OID 16508)
-- Name: ticket specified_location_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket ALTER COLUMN specified_location_id SET DEFAULT nextval('public.ticket_specified_location_id_seq'::regclass);


--
-- TOC entry 3369 (class 2604 OID 16541)
-- Name: track_status_ticket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket ALTER COLUMN id SET DEFAULT nextval('public.track_status_ticket_id_seq'::regclass);


--
-- TOC entry 3370 (class 2604 OID 16542)
-- Name: track_status_ticket ticket_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket ALTER COLUMN ticket_id SET DEFAULT nextval('public.track_status_ticket_ticket_id_seq'::regclass);


--
-- TOC entry 3371 (class 2604 OID 16543)
-- Name: track_status_ticket status_ticket_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket ALTER COLUMN status_ticket_id SET DEFAULT nextval('public.track_status_ticket_status_ticket_id_seq'::regclass);


--
-- TOC entry 3592 (class 0 OID 16420)
-- Dependencies: 223
-- Data for Name: area; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.area (id, department_id, name, is_active, created_at, updated_at) FROM stdin;
1	1	Building	t	2025-07-12 07:32:05.979478+00	2025-07-12 07:32:05.979478+00
2	1	Electrical	t	2025-07-12 07:32:05.981947+00	2025-07-12 07:32:05.981947+00
3	1	Office	t	2025-07-12 07:32:05.982805+00	2025-07-12 07:32:05.982805+00
4	2	Maintenance 1	t	2025-07-12 07:32:05.983605+00	2025-07-12 07:32:05.983605+00
5	2	Maintenance 2	t	2025-07-12 07:32:05.984402+00	2025-07-12 07:32:05.984402+00
6	2	Maintenance Support	t	2025-07-12 07:32:05.985197+00	2025-07-12 07:32:05.985197+00
7	3	Pengukuran	t	2025-07-12 07:32:05.985945+00	2025-07-12 07:32:05.985945+00
8	3	Pengujian	t	2025-07-12 07:32:05.986603+00	2025-07-12 07:32:05.986603+00
\.


--
-- TOC entry 3589 (class 0 OID 16404)
-- Dependencies: 220
-- Data for Name: department; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.department (id, name, receive_job, is_active, created_at, updated_at) FROM stdin;
1	HRGA	t	t	2025-07-12 07:32:05.972199+00	2025-07-12 07:32:05.972199+00
2	Maintenance	t	t	2025-07-12 07:32:05.973629+00	2025-07-12 07:32:05.973629+00
3	Quality	t	t	2025-07-12 07:32:05.97451+00	2025-07-12 07:32:05.97451+00
4	PE	t	t	2025-07-12 07:32:05.975168+00	2025-07-12 07:32:05.975168+00
5	Office	f	t	2025-07-12 07:32:05.975864+00	2025-07-12 07:32:05.975864+00
6	Marketing	f	t	2025-07-12 07:32:05.976535+00	2025-07-12 07:32:05.976535+00
7	Finance	f	t	2025-07-12 07:32:05.9772+00	2025-07-12 07:32:05.9772+00
8	Operation	f	t	2025-07-12 07:32:05.977841+00	2025-07-12 07:32:05.977841+00
\.


--
-- TOC entry 3595 (class 0 OID 16441)
-- Dependencies: 226
-- Data for Name: employee; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.employee (npk, department_id, area_id, name, "position", is_active, created_at, updated_at) FROM stdin;
EMP0001	1	1	Nadia Kusumo	Head of Department	t	2025-07-12 07:32:05.999088+00	2025-07-12 07:32:05.999088+00
EMP0002	1	2	Sari Nugroho	Section	t	2025-07-12 07:32:06.000282+00	2025-07-12 07:32:06.000282+00
EMP0003	1	2	Hadi Kusumo	Section	t	2025-07-12 07:32:06.001285+00	2025-07-12 07:32:06.001285+00
EMP0004	1	1	Lia Susanto	Leader	t	2025-07-12 07:32:06.002186+00	2025-07-12 07:32:06.002186+00
EMP0005	1	2	Lia Setiawan	Leader	t	2025-07-12 07:32:06.002947+00	2025-07-12 07:32:06.002947+00
EMP0006	1	1	Tono Setiawan	Staff	t	2025-07-12 07:32:06.003687+00	2025-07-12 07:32:06.003687+00
EMP0007	1	2	Fajar Susanto	Staff	t	2025-07-12 07:32:06.004394+00	2025-07-12 07:32:06.004394+00
EMP0008	1	3	Eka Wijaya	Staff	t	2025-07-12 07:32:06.004988+00	2025-07-12 07:32:06.004988+00
EMP0009	1	2	Lia Purnama	Staff	t	2025-07-12 07:32:06.00558+00	2025-07-12 07:32:06.00558+00
EMP0010	1	3	Sari Pratama	Staff	t	2025-07-12 07:32:06.006134+00	2025-07-12 07:32:06.006134+00
EMP0011	2	6	Deni Lestari	Head of Department	t	2025-07-12 07:32:06.006785+00	2025-07-12 07:32:06.006785+00
EMP0012	2	4	Deni Lestari	Section	t	2025-07-12 07:32:06.007391+00	2025-07-12 07:32:06.007391+00
EMP0013	2	6	Kartika Pratama	Section	t	2025-07-12 07:32:06.007975+00	2025-07-12 07:32:06.007975+00
EMP0014	2	6	Deni Kusumo	Leader	t	2025-07-12 07:32:06.008638+00	2025-07-12 07:32:06.008638+00
EMP0015	2	5	Fajar Nugroho	Leader	t	2025-07-12 07:32:06.009281+00	2025-07-12 07:32:06.009281+00
EMP0016	2	4	Sari Kusumo	Staff	t	2025-07-12 07:32:06.009814+00	2025-07-12 07:32:06.009814+00
EMP0017	2	5	Hadi Pratama	Staff	t	2025-07-12 07:32:06.010384+00	2025-07-12 07:32:06.010384+00
EMP0018	2	4	Deni Setiawan	Staff	t	2025-07-12 07:32:06.01092+00	2025-07-12 07:32:06.01092+00
EMP0019	2	5	Lia Hidayat	Staff	t	2025-07-12 07:32:06.01163+00	2025-07-12 07:32:06.01163+00
EMP0020	2	5	Gita Susanto	Staff	t	2025-07-12 07:32:06.012384+00	2025-07-12 07:32:06.012384+00
EMP0021	2	5	Deni Lestari	Staff	t	2025-07-12 07:32:06.013243+00	2025-07-12 07:32:06.013243+00
EMP0022	3	8	Fajar Kusumo	Head of Department	t	2025-07-12 07:32:06.014049+00	2025-07-12 07:32:06.014049+00
EMP0023	3	7	Fajar Setiawan	Section	t	2025-07-12 07:32:06.014885+00	2025-07-12 07:32:06.014885+00
EMP0024	3	7	Eka Hidayat	Section	t	2025-07-12 07:32:06.015482+00	2025-07-12 07:32:06.015482+00
EMP0025	3	8	Sari Nugroho	Leader	t	2025-07-12 07:32:06.016068+00	2025-07-12 07:32:06.016068+00
EMP0026	3	8	Oscar Wahyuni	Leader	t	2025-07-12 07:32:06.016815+00	2025-07-12 07:32:06.016815+00
EMP0027	3	7	Deni Kusumo	Leader	t	2025-07-12 07:32:06.017372+00	2025-07-12 07:32:06.017372+00
EMP0028	3	7	Cahyo Pratama	Staff	t	2025-07-12 07:32:06.017938+00	2025-07-12 07:32:06.017938+00
EMP0029	3	8	Adi Kusumo	Staff	t	2025-07-12 07:32:06.01876+00	2025-07-12 07:32:06.01876+00
EMP0030	3	7	Nadia Lestari	Staff	t	2025-07-12 07:32:06.019311+00	2025-07-12 07:32:06.019311+00
EMP0031	3	7	Cahyo Setiawan	Staff	t	2025-07-12 07:32:06.02003+00	2025-07-12 07:32:06.02003+00
EMP0032	3	7	Rina Lestari	Staff	t	2025-07-12 07:32:06.020581+00	2025-07-12 07:32:06.020581+00
EMP0033	3	7	Indra Nugroho	Staff	t	2025-07-12 07:32:06.021148+00	2025-07-12 07:32:06.021148+00
EMP0034	3	7	Rina Lestari	Staff	t	2025-07-12 07:32:06.021845+00	2025-07-12 07:32:06.021845+00
EMP0035	3	7	Hadi Susanto	Staff	t	2025-07-12 07:32:06.02239+00	2025-07-12 07:32:06.02239+00
EMP0036	3	8	Putra Setiawan	Staff	t	2025-07-12 07:32:06.023085+00	2025-07-12 07:32:06.023085+00
EMP0037	4	\N	Tono Susanto	Head of Department	t	2025-07-12 07:32:06.023825+00	2025-07-12 07:32:06.023825+00
EMP0038	4	\N	Wati Hidayat	Section	t	2025-07-12 07:32:06.024548+00	2025-07-12 07:32:06.024548+00
EMP0039	4	\N	Adi Lestari	Section	t	2025-07-12 07:32:06.025233+00	2025-07-12 07:32:06.025233+00
EMP0040	4	\N	Hadi Setiawan	Staff	t	2025-07-12 07:32:06.0258+00	2025-07-12 07:32:06.0258+00
EMP0041	4	\N	Oscar Wahyuni	Staff	t	2025-07-12 07:32:06.02636+00	2025-07-12 07:32:06.02636+00
EMP0042	4	\N	Budi Pratama	Staff	t	2025-07-12 07:32:06.026915+00	2025-07-12 07:32:06.026915+00
EMP0043	4	\N	Eka Susanto	Staff	t	2025-07-12 07:32:06.027458+00	2025-07-12 07:32:06.027458+00
EMP0044	4	\N	Wati Pratama	Staff	t	2025-07-12 07:32:06.028256+00	2025-07-12 07:32:06.028256+00
EMP0045	4	\N	Nadia Hidayat	Staff	t	2025-07-12 07:32:06.029021+00	2025-07-12 07:32:06.029021+00
EMP0046	4	\N	Nadia Kusumo	Staff	t	2025-07-12 07:32:06.029676+00	2025-07-12 07:32:06.029676+00
EMP0047	5	\N	Mega Wijaya	Head of Department	t	2025-07-12 07:32:06.03033+00	2025-07-12 07:32:06.03033+00
EMP0048	5	\N	Nadia Nugroho	Section	t	2025-07-12 07:32:06.030927+00	2025-07-12 07:32:06.030927+00
EMP0049	5	\N	Fajar Nugroho	Section	t	2025-07-12 07:32:06.031486+00	2025-07-12 07:32:06.031486+00
EMP0050	5	\N	Fajar Purnama	Staff	t	2025-07-12 07:32:06.032057+00	2025-07-12 07:32:06.032057+00
EMP0051	5	\N	Nadia Purnama	Staff	t	2025-07-12 07:32:06.032621+00	2025-07-12 07:32:06.032621+00
EMP0052	5	\N	Hadi Lestari	Staff	t	2025-07-12 07:32:06.03318+00	2025-07-12 07:32:06.03318+00
EMP0053	5	\N	Gita Setiawan	Staff	t	2025-07-12 07:32:06.033736+00	2025-07-12 07:32:06.033736+00
EMP0054	5	\N	Sari Hidayat	Staff	t	2025-07-12 07:32:06.034289+00	2025-07-12 07:32:06.034289+00
EMP0055	6	\N	Joko Kusumo	Head of Department	t	2025-07-12 07:32:06.034901+00	2025-07-12 07:32:06.034901+00
EMP0056	6	\N	Indra Susanto	Section	t	2025-07-12 07:32:06.035468+00	2025-07-12 07:32:06.035468+00
EMP0057	6	\N	Fajar Nugroho	Section	t	2025-07-12 07:32:06.036029+00	2025-07-12 07:32:06.036029+00
EMP0058	6	\N	Joko Setiawan	Staff	t	2025-07-12 07:32:06.036584+00	2025-07-12 07:32:06.036584+00
EMP0059	6	\N	Deni Wahyuni	Staff	t	2025-07-12 07:32:06.037262+00	2025-07-12 07:32:06.037262+00
EMP0060	6	\N	Sari Purnama	Staff	t	2025-07-12 07:32:06.037875+00	2025-07-12 07:32:06.037875+00
EMP0061	6	\N	Fajar Wijaya	Staff	t	2025-07-12 07:32:06.038453+00	2025-07-12 07:32:06.038453+00
EMP0062	6	\N	Lia Lestari	Staff	t	2025-07-12 07:32:06.039046+00	2025-07-12 07:32:06.039046+00
EMP0063	6	\N	Sari Hidayat	Staff	t	2025-07-12 07:32:06.039663+00	2025-07-12 07:32:06.039663+00
EMP0064	7	\N	Gita Wahyuni	Head of Department	t	2025-07-12 07:32:06.040317+00	2025-07-12 07:32:06.040317+00
EMP0065	7	\N	Budi Lestari	Section	t	2025-07-12 07:32:06.040933+00	2025-07-12 07:32:06.040933+00
EMP0066	7	\N	Eka Pratama	Section	t	2025-07-12 07:32:06.041535+00	2025-07-12 07:32:06.041535+00
EMP0067	7	\N	Cahyo Susanto	Staff	t	2025-07-12 07:32:06.042129+00	2025-07-12 07:32:06.042129+00
EMP0068	7	\N	Budi Setiawan	Staff	t	2025-07-12 07:32:06.042705+00	2025-07-12 07:32:06.042705+00
EMP0069	7	\N	Adi Susanto	Staff	t	2025-07-12 07:32:06.043282+00	2025-07-12 07:32:06.043282+00
EMP0070	7	\N	Eka Hidayat	Staff	t	2025-07-12 07:32:06.043912+00	2025-07-12 07:32:06.043912+00
EMP0071	7	\N	Nadia Lestari	Staff	t	2025-07-12 07:32:06.044527+00	2025-07-12 07:32:06.044527+00
EMP0072	7	\N	Gita Kusumo	Staff	t	2025-07-12 07:32:06.045181+00	2025-07-12 07:32:06.045181+00
EMP0073	7	\N	Cahyo Lestari	Staff	t	2025-07-12 07:32:06.045804+00	2025-07-12 07:32:06.045804+00
EMP0074	8	\N	Sari Purnama	Head of Department	t	2025-07-12 07:32:06.046456+00	2025-07-12 07:32:06.046456+00
EMP0075	8	\N	Gita Susanto	Section	t	2025-07-12 07:32:06.047062+00	2025-07-12 07:32:06.047062+00
EMP0076	8	\N	Gita Hidayat	Section	t	2025-07-12 07:32:06.047617+00	2025-07-12 07:32:06.047617+00
EMP0077	8	\N	Wati Kusumo	Staff	t	2025-07-12 07:32:06.048174+00	2025-07-12 07:32:06.048174+00
EMP0078	8	\N	Lia Susanto	Staff	t	2025-07-12 07:32:06.048731+00	2025-07-12 07:32:06.048731+00
EMP0079	8	\N	Gita Hidayat	Staff	t	2025-07-12 07:32:06.049295+00	2025-07-12 07:32:06.049295+00
EMP0080	8	\N	Deni Pratama	Staff	t	2025-07-12 07:32:06.049854+00	2025-07-12 07:32:06.049854+00
EMP0081	8	\N	Wati Wijaya	Staff	t	2025-07-12 07:32:06.050419+00	2025-07-12 07:32:06.050419+00
\.


--
-- TOC entry 3612 (class 0 OID 16562)
-- Dependencies: 243
-- Data for Name: job; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.job (id, ticket_id, pic_job, job_priority, report_file, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 3597 (class 0 OID 16464)
-- Dependencies: 228
-- Data for Name: physical_location; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.physical_location (id, name, is_active, created_at, updated_at) FROM stdin;
1	Forging	t	2025-07-12 07:32:05.988023+00	2025-07-12 07:32:05.988023+00
2	Production	t	2025-07-12 07:32:05.989243+00	2025-07-12 07:32:05.989243+00
3	Log	t	2025-07-12 07:32:05.990076+00	2025-07-12 07:32:05.990076+00
4	Building Office	f	2025-07-12 07:32:05.990723+00	2025-07-12 07:32:05.990723+00
\.


--
-- TOC entry 3615 (class 0 OID 16587)
-- Dependencies: 246
-- Data for Name: rejected_ticket; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rejected_ticket (id, ticket_id, rejector, feedback, already_seen, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 3600 (class 0 OID 16479)
-- Dependencies: 231
-- Data for Name: specified_location; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.specified_location (id, physical_location_id, name, is_active, created_at, updated_at) FROM stdin;
1	1	F 1	t	2025-07-12 07:32:05.99192+00	2025-07-12 07:32:05.99192+00
2	1	F 2	t	2025-07-12 07:32:05.992956+00	2025-07-12 07:32:05.992956+00
3	1	F 3	t	2025-07-12 07:32:05.993695+00	2025-07-12 07:32:05.993695+00
4	1	F 4	t	2025-07-12 07:32:05.994364+00	2025-07-12 07:32:05.994364+00
5	2	Machine Production	t	2025-07-12 07:32:05.99511+00	2025-07-12 07:32:05.99511+00
6	2	Tool Production	t	2025-07-12 07:32:05.995782+00	2025-07-12 07:32:05.995782+00
7	2	Support Production	t	2025-07-12 07:32:05.996438+00	2025-07-12 07:32:05.996438+00
8	3	Input Log	t	2025-07-12 07:32:05.997145+00	2025-07-12 07:32:05.997145+00
9	3	Output Log	t	2025-07-12 07:32:05.997748+00	2025-07-12 07:32:05.997748+00
\.


--
-- TOC entry 3587 (class 0 OID 16390)
-- Dependencies: 218
-- Data for Name: status_ticket; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.status_ticket (id, name, sequence, is_active, created_at, updated_at) FROM stdin;
1	Dibatalkan	-100	t	2025-07-12 07:32:05.963309+00	2025-07-12 07:32:05.963309+00
2	Approval Section	-2	t	2025-07-12 07:32:05.965858+00	2025-07-12 07:32:05.965858+00
3	Approval Department	-1	t	2025-07-12 07:32:05.967244+00	2025-07-12 07:32:05.967244+00
4	Menunggu Job	0	t	2025-07-12 07:32:05.968298+00	2025-07-12 07:32:05.968298+00
5	Dikerjakan	1	t	2025-07-12 07:32:05.96903+00	2025-07-12 07:32:05.96903+00
6	Job Selesai	2	t	2025-07-12 07:32:05.969685+00	2025-07-12 07:32:05.969685+00
7	Tiket selesai	3	t	2025-07-12 07:32:05.970524+00	2025-07-12 07:32:05.970524+00
\.


--
-- TOC entry 3605 (class 0 OID 16502)
-- Dependencies: 236
-- Data for Name: ticket; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.ticket (id, requestor, department_target_id, physical_location_id, specified_location_id, description, ticket_priority, support_file, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 3609 (class 0 OID 16538)
-- Dependencies: 240
-- Data for Name: track_status_ticket; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.track_status_ticket (id, ticket_id, status_ticket_id, start_date, finish_date) FROM stdin;
\.


--
-- TOC entry 3641 (class 0 OID 0)
-- Dependencies: 222
-- Name: area_department_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.area_department_id_seq', 1, false);


--
-- TOC entry 3642 (class 0 OID 0)
-- Dependencies: 221
-- Name: area_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.area_id_seq', 8, true);


--
-- TOC entry 3643 (class 0 OID 0)
-- Dependencies: 219
-- Name: department_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.department_id_seq', 8, true);


--
-- TOC entry 3644 (class 0 OID 0)
-- Dependencies: 225
-- Name: employee_area_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.employee_area_id_seq', 1, false);


--
-- TOC entry 3645 (class 0 OID 0)
-- Dependencies: 224
-- Name: employee_department_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.employee_department_id_seq', 1, false);


--
-- TOC entry 3646 (class 0 OID 0)
-- Dependencies: 241
-- Name: job_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.job_id_seq', 1, false);


--
-- TOC entry 3647 (class 0 OID 0)
-- Dependencies: 242
-- Name: job_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.job_ticket_id_seq', 1, false);


--
-- TOC entry 3648 (class 0 OID 0)
-- Dependencies: 227
-- Name: physical_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.physical_location_id_seq', 4, true);


--
-- TOC entry 3649 (class 0 OID 0)
-- Dependencies: 244
-- Name: rejected_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rejected_ticket_id_seq', 1, false);


--
-- TOC entry 3650 (class 0 OID 0)
-- Dependencies: 245
-- Name: rejected_ticket_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rejected_ticket_ticket_id_seq', 1, false);


--
-- TOC entry 3651 (class 0 OID 0)
-- Dependencies: 229
-- Name: specified_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.specified_location_id_seq', 9, true);


--
-- TOC entry 3652 (class 0 OID 0)
-- Dependencies: 230
-- Name: specified_location_physical_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.specified_location_physical_location_id_seq', 1, false);


--
-- TOC entry 3653 (class 0 OID 0)
-- Dependencies: 217
-- Name: status_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.status_ticket_id_seq', 7, true);


--
-- TOC entry 3654 (class 0 OID 0)
-- Dependencies: 233
-- Name: ticket_department_target_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ticket_department_target_id_seq', 1, false);


--
-- TOC entry 3655 (class 0 OID 0)
-- Dependencies: 232
-- Name: ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ticket_id_seq', 1, false);


--
-- TOC entry 3656 (class 0 OID 0)
-- Dependencies: 234
-- Name: ticket_physical_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ticket_physical_location_id_seq', 1, false);


--
-- TOC entry 3657 (class 0 OID 0)
-- Dependencies: 235
-- Name: ticket_specified_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ticket_specified_location_id_seq', 1, false);


--
-- TOC entry 3658 (class 0 OID 0)
-- Dependencies: 237
-- Name: track_status_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.track_status_ticket_id_seq', 1, false);


--
-- TOC entry 3659 (class 0 OID 0)
-- Dependencies: 239
-- Name: track_status_ticket_status_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.track_status_ticket_status_ticket_id_seq', 1, false);


--
-- TOC entry 3660 (class 0 OID 0)
-- Dependencies: 238
-- Name: track_status_ticket_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.track_status_ticket_ticket_id_seq', 1, false);


--
-- TOC entry 3391 (class 2606 OID 16431)
-- Name: area area_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.area
    ADD CONSTRAINT area_pkey PRIMARY KEY (id);


--
-- TOC entry 3393 (class 2606 OID 16433)
-- Name: area area_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.area
    ADD CONSTRAINT area_unique UNIQUE (name, department_id);


--
-- TOC entry 3387 (class 2606 OID 16415)
-- Name: department department_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.department
    ADD CONSTRAINT department_pkey PRIMARY KEY (id);


--
-- TOC entry 3389 (class 2606 OID 16417)
-- Name: department department_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.department
    ADD CONSTRAINT department_unique UNIQUE (name);


--
-- TOC entry 3395 (class 2606 OID 16452)
-- Name: employee employee_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT employee_pkey PRIMARY KEY (npk);


--
-- TOC entry 3411 (class 2606 OID 16572)
-- Name: job job_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT job_pkey PRIMARY KEY (id);


--
-- TOC entry 3413 (class 2606 OID 16574)
-- Name: job job_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT job_unique UNIQUE (ticket_id);


--
-- TOC entry 3397 (class 2606 OID 16474)
-- Name: physical_location physical_location_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.physical_location
    ADD CONSTRAINT physical_location_pkey PRIMARY KEY (id);


--
-- TOC entry 3399 (class 2606 OID 16476)
-- Name: physical_location physical_location_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.physical_location
    ADD CONSTRAINT physical_location_unique UNIQUE (name);


--
-- TOC entry 3415 (class 2606 OID 16598)
-- Name: rejected_ticket rejected_ticket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rejected_ticket
    ADD CONSTRAINT rejected_ticket_pkey PRIMARY KEY (id);


--
-- TOC entry 3401 (class 2606 OID 16490)
-- Name: specified_location specified_location_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.specified_location
    ADD CONSTRAINT specified_location_pkey PRIMARY KEY (id);


--
-- TOC entry 3403 (class 2606 OID 16492)
-- Name: specified_location specified_location_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.specified_location
    ADD CONSTRAINT specified_location_unique UNIQUE (physical_location_id, name);


--
-- TOC entry 3383 (class 2606 OID 16400)
-- Name: status_ticket status_ticket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.status_ticket
    ADD CONSTRAINT status_ticket_pkey PRIMARY KEY (id);


--
-- TOC entry 3385 (class 2606 OID 16402)
-- Name: status_ticket status_ticket_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.status_ticket
    ADD CONSTRAINT status_ticket_unique UNIQUE (name, sequence);


--
-- TOC entry 3405 (class 2606 OID 16514)
-- Name: ticket ticket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket
    ADD CONSTRAINT ticket_pkey PRIMARY KEY (id);


--
-- TOC entry 3407 (class 2606 OID 16548)
-- Name: track_status_ticket track_status_ticket_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket
    ADD CONSTRAINT track_status_ticket_id UNIQUE (ticket_id, status_ticket_id);


--
-- TOC entry 3409 (class 2606 OID 16546)
-- Name: track_status_ticket track_status_ticket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket
    ADD CONSTRAINT track_status_ticket_pkey PRIMARY KEY (id);


--
-- TOC entry 3432 (class 2620 OID 16610)
-- Name: area area_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER area_set_timestamp BEFORE UPDATE ON public.area FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3431 (class 2620 OID 16611)
-- Name: department department_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER department_set_timestamp BEFORE UPDATE ON public.department FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3433 (class 2620 OID 16612)
-- Name: employee employee; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER employee BEFORE UPDATE ON public.employee FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3437 (class 2620 OID 16613)
-- Name: job job_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER job_set_timestamp BEFORE UPDATE ON public.job FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3434 (class 2620 OID 16614)
-- Name: physical_location physical_location_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER physical_location_set_timestamp BEFORE UPDATE ON public.physical_location FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3438 (class 2620 OID 16615)
-- Name: rejected_ticket rejected_ticket_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER rejected_ticket_set_timestamp BEFORE UPDATE ON public.rejected_ticket FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3435 (class 2620 OID 16616)
-- Name: specified_location specified_location_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER specified_location_set_timestamp BEFORE UPDATE ON public.specified_location FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3430 (class 2620 OID 16617)
-- Name: status_ticket status_ticket_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER status_ticket_set_timestamp BEFORE UPDATE ON public.status_ticket FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3436 (class 2620 OID 16618)
-- Name: ticket ticket_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER ticket_set_timestamp BEFORE UPDATE ON public.ticket FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3417 (class 2606 OID 16458)
-- Name: employee area_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT area_id FOREIGN KEY (area_id) REFERENCES public.area(id);


--
-- TOC entry 3416 (class 2606 OID 16434)
-- Name: area department_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.area
    ADD CONSTRAINT department_id FOREIGN KEY (department_id) REFERENCES public.department(id);


--
-- TOC entry 3418 (class 2606 OID 16453)
-- Name: employee department_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT department_id FOREIGN KEY (department_id) REFERENCES public.department(id);


--
-- TOC entry 3420 (class 2606 OID 16520)
-- Name: ticket department_target_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket
    ADD CONSTRAINT department_target_id FOREIGN KEY (department_target_id) REFERENCES public.department(id);


--
-- TOC entry 3419 (class 2606 OID 16493)
-- Name: specified_location physical_location_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.specified_location
    ADD CONSTRAINT physical_location_id FOREIGN KEY (physical_location_id) REFERENCES public.physical_location(id);


--
-- TOC entry 3421 (class 2606 OID 16525)
-- Name: ticket physical_location_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket
    ADD CONSTRAINT physical_location_id FOREIGN KEY (physical_location_id) REFERENCES public.physical_location(id);


--
-- TOC entry 3426 (class 2606 OID 16580)
-- Name: job pic_job; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT pic_job FOREIGN KEY (pic_job) REFERENCES public.employee(npk);


--
-- TOC entry 3428 (class 2606 OID 16604)
-- Name: rejected_ticket rejector; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rejected_ticket
    ADD CONSTRAINT rejector FOREIGN KEY (rejector) REFERENCES public.employee(npk);


--
-- TOC entry 3422 (class 2606 OID 16515)
-- Name: ticket requestor; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket
    ADD CONSTRAINT requestor FOREIGN KEY (requestor) REFERENCES public.employee(npk);


--
-- TOC entry 3423 (class 2606 OID 16530)
-- Name: ticket specified_location_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket
    ADD CONSTRAINT specified_location_id FOREIGN KEY (specified_location_id) REFERENCES public.specified_location(id);


--
-- TOC entry 3424 (class 2606 OID 16554)
-- Name: track_status_ticket status_ticket_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket
    ADD CONSTRAINT status_ticket_id FOREIGN KEY (status_ticket_id) REFERENCES public.status_ticket(id);


--
-- TOC entry 3425 (class 2606 OID 16549)
-- Name: track_status_ticket ticket_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket
    ADD CONSTRAINT ticket_id FOREIGN KEY (ticket_id) REFERENCES public.ticket(id);


--
-- TOC entry 3427 (class 2606 OID 16575)
-- Name: job ticket_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT ticket_id FOREIGN KEY (ticket_id) REFERENCES public.ticket(id);


--
-- TOC entry 3429 (class 2606 OID 16599)
-- Name: rejected_ticket ticket_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rejected_ticket
    ADD CONSTRAINT ticket_id FOREIGN KEY (ticket_id) REFERENCES public.ticket(id);


-- Completed on 2025-07-12 14:44:35

--
-- PostgreSQL database dump complete
--

