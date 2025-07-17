--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5
-- Dumped by pg_dump version 17.5

-- Started on 2025-07-16 11:45:52

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
-- TOC entry 270 (class 1255 OID 16609)
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
-- TOC entry 3747 (class 0 OID 0)
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
-- TOC entry 3748 (class 0 OID 0)
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
-- TOC entry 3749 (class 0 OID 0)
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
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    position_id smallint NOT NULL
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
-- TOC entry 3750 (class 0 OID 0)
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
-- TOC entry 3751 (class 0 OID 0)
-- Dependencies: 224
-- Name: employee_department_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.employee_department_id_seq OWNED BY public.employee.department_id;


--
-- TOC entry 253 (class 1259 OID 18261)
-- Name: employee_position_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.employee_position_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.employee_position_id_seq OWNER TO postgres;

--
-- TOC entry 3752 (class 0 OID 0)
-- Dependencies: 253
-- Name: employee_position_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.employee_position_id_seq OWNED BY public.employee.position_id;


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
-- TOC entry 3753 (class 0 OID 0)
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
-- TOC entry 3754 (class 0 OID 0)
-- Dependencies: 242
-- Name: job_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.job_ticket_id_seq OWNED BY public.job.ticket_id;


--
-- TOC entry 265 (class 1259 OID 20527)
-- Name: permission; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.permission (
    id smallint NOT NULL,
    name text NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.permission OWNER TO postgres;

--
-- TOC entry 264 (class 1259 OID 20526)
-- Name: permission_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.permission_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.permission_id_seq OWNER TO postgres;

--
-- TOC entry 3755 (class 0 OID 0)
-- Dependencies: 264
-- Name: permission_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.permission_id_seq OWNED BY public.permission.id;


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
-- TOC entry 3756 (class 0 OID 0)
-- Dependencies: 227
-- Name: physical_location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.physical_location_id_seq OWNED BY public.physical_location.id;


--
-- TOC entry 252 (class 1259 OID 18249)
-- Name: position; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."position" (
    id smallint NOT NULL,
    name text NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public."position" OWNER TO postgres;

--
-- TOC entry 251 (class 1259 OID 18248)
-- Name: position_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.position_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.position_id_seq OWNER TO postgres;

--
-- TOC entry 3757 (class 0 OID 0)
-- Dependencies: 251
-- Name: position_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.position_id_seq OWNED BY public."position".id;


--
-- TOC entry 268 (class 1259 OID 20543)
-- Name: position_permission; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.position_permission (
    position_id smallint NOT NULL,
    permission_id smallint NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.position_permission OWNER TO postgres;

--
-- TOC entry 267 (class 1259 OID 20542)
-- Name: position_permission_permission_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.position_permission_permission_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.position_permission_permission_id_seq OWNER TO postgres;

--
-- TOC entry 3758 (class 0 OID 0)
-- Dependencies: 267
-- Name: position_permission_permission_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.position_permission_permission_id_seq OWNED BY public.position_permission.permission_id;


--
-- TOC entry 266 (class 1259 OID 20541)
-- Name: position_permission_position_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.position_permission_position_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.position_permission_position_id_seq OWNER TO postgres;

--
-- TOC entry 3759 (class 0 OID 0)
-- Dependencies: 266
-- Name: position_permission_position_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.position_permission_position_id_seq OWNED BY public.position_permission.position_id;


--
-- TOC entry 263 (class 1259 OID 18324)
-- Name: position_to_workflow_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.position_to_workflow_mapping (
    id smallint NOT NULL,
    position_id smallint NOT NULL,
    workflow_id smallint NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.position_to_workflow_mapping OWNER TO postgres;

--
-- TOC entry 260 (class 1259 OID 18321)
-- Name: position_to_workflow_mapping_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.position_to_workflow_mapping_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.position_to_workflow_mapping_id_seq OWNER TO postgres;

--
-- TOC entry 3760 (class 0 OID 0)
-- Dependencies: 260
-- Name: position_to_workflow_mapping_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.position_to_workflow_mapping_id_seq OWNED BY public.position_to_workflow_mapping.id;


--
-- TOC entry 261 (class 1259 OID 18322)
-- Name: position_to_workflow_mapping_position_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.position_to_workflow_mapping_position_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.position_to_workflow_mapping_position_id_seq OWNER TO postgres;

--
-- TOC entry 3761 (class 0 OID 0)
-- Dependencies: 261
-- Name: position_to_workflow_mapping_position_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.position_to_workflow_mapping_position_id_seq OWNED BY public.position_to_workflow_mapping.position_id;


--
-- TOC entry 262 (class 1259 OID 18323)
-- Name: position_to_workflow_mapping_workflow_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.position_to_workflow_mapping_workflow_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.position_to_workflow_mapping_workflow_id_seq OWNER TO postgres;

--
-- TOC entry 3762 (class 0 OID 0)
-- Dependencies: 262
-- Name: position_to_workflow_mapping_workflow_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.position_to_workflow_mapping_workflow_id_seq OWNED BY public.position_to_workflow_mapping.workflow_id;


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
-- TOC entry 3763 (class 0 OID 0)
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
-- TOC entry 3764 (class 0 OID 0)
-- Dependencies: 245
-- Name: rejected_ticket_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.rejected_ticket_ticket_id_seq OWNED BY public.rejected_ticket.ticket_id;


--
-- TOC entry 249 (class 1259 OID 18206)
-- Name: section_status_ticket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.section_status_ticket (
    id smallint NOT NULL,
    name text NOT NULL,
    sequence smallint NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.section_status_ticket OWNER TO postgres;

--
-- TOC entry 248 (class 1259 OID 18205)
-- Name: section_status_ticket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.section_status_ticket_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.section_status_ticket_id_seq OWNER TO postgres;

--
-- TOC entry 3765 (class 0 OID 0)
-- Dependencies: 248
-- Name: section_status_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.section_status_ticket_id_seq OWNED BY public.section_status_ticket.id;


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
-- TOC entry 3766 (class 0 OID 0)
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
-- TOC entry 3767 (class 0 OID 0)
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
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    section_id smallint NOT NULL
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
-- TOC entry 3768 (class 0 OID 0)
-- Dependencies: 217
-- Name: status_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.status_ticket_id_seq OWNED BY public.status_ticket.id;


--
-- TOC entry 250 (class 1259 OID 18227)
-- Name: status_ticket_section_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.status_ticket_section_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.status_ticket_section_id_seq OWNER TO postgres;

--
-- TOC entry 3769 (class 0 OID 0)
-- Dependencies: 250
-- Name: status_ticket_section_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.status_ticket_section_id_seq OWNED BY public.status_ticket.section_id;


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
-- TOC entry 3770 (class 0 OID 0)
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
-- TOC entry 3771 (class 0 OID 0)
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
-- TOC entry 3772 (class 0 OID 0)
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
-- TOC entry 3773 (class 0 OID 0)
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
-- TOC entry 3774 (class 0 OID 0)
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
-- TOC entry 3775 (class 0 OID 0)
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
-- TOC entry 3776 (class 0 OID 0)
-- Dependencies: 238
-- Name: track_status_ticket_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.track_status_ticket_ticket_id_seq OWNED BY public.track_status_ticket.ticket_id;


--
-- TOC entry 247 (class 1259 OID 16624)
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
-- TOC entry 269 (class 1259 OID 20564)
-- Name: view_ticket_list; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.view_ticket_list AS
 WITH ticket_actual_start AS (
         SELECT tst.ticket_id,
            min(tst.start_date) AS actual_start_date
           FROM ((public.track_status_ticket tst
             JOIN public.status_ticket st ON ((tst.status_ticket_id = st.id)))
             JOIN public.section_status_ticket sst ON ((st.section_id = sst.id)))
          WHERE (sst.name = 'Actual Section'::text)
          GROUP BY tst.ticket_id
        )
 SELECT t.id AS ticket_id,
    t.description,
    t.department_target_id,
    t.ticket_priority,
    j.job_priority,
        CASE
            WHEN (tas.actual_start_date IS NOT NULL) THEN ((now())::date - (tas.actual_start_date)::date)
            ELSE NULL::integer
        END AS ticket_age_days,
    req_emp.name AS requestor_name,
    req_dept.name AS requestor_department,
    ( SELECT st.name
           FROM (public.track_status_ticket tst
             JOIN public.status_ticket st ON ((tst.status_ticket_id = st.id)))
          WHERE (tst.ticket_id = t.id)
          ORDER BY tst.start_date DESC, tst.id DESC
         LIMIT 1) AS current_status,
    pic_emp.name AS pic_name,
    pic_area.name AS pic_area_name,
    phys_loc.name AS location_name,
    spec_loc.name AS specified_location_name
   FROM ((((((((public.ticket t
     LEFT JOIN public.job j ON ((t.id = j.ticket_id)))
     LEFT JOIN ticket_actual_start tas ON ((t.id = tas.ticket_id)))
     JOIN public.employee req_emp ON ((t.requestor = req_emp.npk)))
     LEFT JOIN public.department req_dept ON ((req_emp.department_id = req_dept.id)))
     LEFT JOIN public.employee pic_emp ON ((j.pic_job = pic_emp.npk)))
     LEFT JOIN public.area pic_area ON ((pic_emp.area_id = pic_area.id)))
     LEFT JOIN public.physical_location phys_loc ON ((t.physical_location_id = phys_loc.id)))
     LEFT JOIN public.specified_location spec_loc ON ((t.specified_location_id = spec_loc.id)));


ALTER VIEW public.view_ticket_list OWNER TO postgres;

--
-- TOC entry 255 (class 1259 OID 18277)
-- Name: workflow; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.workflow (
    id smallint NOT NULL,
    name text NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.workflow OWNER TO postgres;

--
-- TOC entry 254 (class 1259 OID 18276)
-- Name: workflow_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.workflow_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.workflow_id_seq OWNER TO postgres;

--
-- TOC entry 3777 (class 0 OID 0)
-- Dependencies: 254
-- Name: workflow_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.workflow_id_seq OWNED BY public.workflow.id;


--
-- TOC entry 259 (class 1259 OID 18294)
-- Name: workflow_step; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.workflow_step (
    id smallint NOT NULL,
    workflow_id smallint NOT NULL,
    status_ticket_id smallint NOT NULL,
    step_sequence smallint NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.workflow_step OWNER TO postgres;

--
-- TOC entry 256 (class 1259 OID 18291)
-- Name: workflow_step_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.workflow_step_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.workflow_step_id_seq OWNER TO postgres;

--
-- TOC entry 3778 (class 0 OID 0)
-- Dependencies: 256
-- Name: workflow_step_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.workflow_step_id_seq OWNED BY public.workflow_step.id;


--
-- TOC entry 258 (class 1259 OID 18293)
-- Name: workflow_step_status_ticket_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.workflow_step_status_ticket_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.workflow_step_status_ticket_id_seq OWNER TO postgres;

--
-- TOC entry 3779 (class 0 OID 0)
-- Dependencies: 258
-- Name: workflow_step_status_ticket_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.workflow_step_status_ticket_id_seq OWNED BY public.workflow_step.status_ticket_id;


--
-- TOC entry 257 (class 1259 OID 18292)
-- Name: workflow_step_workflow_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.workflow_step_workflow_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.workflow_step_workflow_id_seq OWNER TO postgres;

--
-- TOC entry 3780 (class 0 OID 0)
-- Dependencies: 257
-- Name: workflow_step_workflow_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.workflow_step_workflow_id_seq OWNED BY public.workflow_step.workflow_id;


--
-- TOC entry 3386 (class 2604 OID 16423)
-- Name: area id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.area ALTER COLUMN id SET DEFAULT nextval('public.area_id_seq'::regclass);


--
-- TOC entry 3381 (class 2604 OID 16407)
-- Name: department id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.department ALTER COLUMN id SET DEFAULT nextval('public.department_id_seq'::regclass);


--
-- TOC entry 3406 (class 2604 OID 16565)
-- Name: job id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job ALTER COLUMN id SET DEFAULT nextval('public.job_id_seq'::regclass);


--
-- TOC entry 3433 (class 2604 OID 20530)
-- Name: permission id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permission ALTER COLUMN id SET DEFAULT nextval('public.permission_id_seq'::regclass);


--
-- TOC entry 3393 (class 2604 OID 16467)
-- Name: physical_location id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.physical_location ALTER COLUMN id SET DEFAULT nextval('public.physical_location_id_seq'::regclass);


--
-- TOC entry 3417 (class 2604 OID 18252)
-- Name: position id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."position" ALTER COLUMN id SET DEFAULT nextval('public.position_id_seq'::regclass);


--
-- TOC entry 3429 (class 2604 OID 18327)
-- Name: position_to_workflow_mapping id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.position_to_workflow_mapping ALTER COLUMN id SET DEFAULT nextval('public.position_to_workflow_mapping_id_seq'::regclass);


--
-- TOC entry 3409 (class 2604 OID 16590)
-- Name: rejected_ticket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rejected_ticket ALTER COLUMN id SET DEFAULT nextval('public.rejected_ticket_id_seq'::regclass);


--
-- TOC entry 3413 (class 2604 OID 18209)
-- Name: section_status_ticket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.section_status_ticket ALTER COLUMN id SET DEFAULT nextval('public.section_status_ticket_id_seq'::regclass);


--
-- TOC entry 3397 (class 2604 OID 16482)
-- Name: specified_location id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.specified_location ALTER COLUMN id SET DEFAULT nextval('public.specified_location_id_seq'::regclass);


--
-- TOC entry 3377 (class 2604 OID 16393)
-- Name: status_ticket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.status_ticket ALTER COLUMN id SET DEFAULT nextval('public.status_ticket_id_seq'::regclass);


--
-- TOC entry 3401 (class 2604 OID 16505)
-- Name: ticket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket ALTER COLUMN id SET DEFAULT nextval('public.ticket_id_seq'::regclass);


--
-- TOC entry 3404 (class 2604 OID 16541)
-- Name: track_status_ticket id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket ALTER COLUMN id SET DEFAULT nextval('public.track_status_ticket_id_seq'::regclass);


--
-- TOC entry 3421 (class 2604 OID 18280)
-- Name: workflow id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workflow ALTER COLUMN id SET DEFAULT nextval('public.workflow_id_seq'::regclass);


--
-- TOC entry 3425 (class 2604 OID 18297)
-- Name: workflow_step id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workflow_step ALTER COLUMN id SET DEFAULT nextval('public.workflow_step_id_seq'::regclass);


--
-- TOC entry 3781 (class 0 OID 0)
-- Dependencies: 222
-- Name: area_department_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.area_department_id_seq', 1, false);


--
-- TOC entry 3782 (class 0 OID 0)
-- Dependencies: 221
-- Name: area_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.area_id_seq', 8, true);


--
-- TOC entry 3783 (class 0 OID 0)
-- Dependencies: 219
-- Name: department_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.department_id_seq', 8, true);


--
-- TOC entry 3784 (class 0 OID 0)
-- Dependencies: 225
-- Name: employee_area_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.employee_area_id_seq', 1, false);


--
-- TOC entry 3785 (class 0 OID 0)
-- Dependencies: 224
-- Name: employee_department_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.employee_department_id_seq', 1, false);


--
-- TOC entry 3786 (class 0 OID 0)
-- Dependencies: 253
-- Name: employee_position_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.employee_position_id_seq', 1, false);


--
-- TOC entry 3787 (class 0 OID 0)
-- Dependencies: 241
-- Name: job_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.job_id_seq', 1, true);


--
-- TOC entry 3788 (class 0 OID 0)
-- Dependencies: 242
-- Name: job_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.job_ticket_id_seq', 1, false);


--
-- TOC entry 3789 (class 0 OID 0)
-- Dependencies: 264
-- Name: permission_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.permission_id_seq', 3, true);


--
-- TOC entry 3790 (class 0 OID 0)
-- Dependencies: 227
-- Name: physical_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.physical_location_id_seq', 4, true);


--
-- TOC entry 3791 (class 0 OID 0)
-- Dependencies: 251
-- Name: position_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.position_id_seq', 4, true);


--
-- TOC entry 3792 (class 0 OID 0)
-- Dependencies: 267
-- Name: position_permission_permission_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.position_permission_permission_id_seq', 1, false);


--
-- TOC entry 3793 (class 0 OID 0)
-- Dependencies: 266
-- Name: position_permission_position_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.position_permission_position_id_seq', 1, false);


--
-- TOC entry 3794 (class 0 OID 0)
-- Dependencies: 260
-- Name: position_to_workflow_mapping_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.position_to_workflow_mapping_id_seq', 4, true);


--
-- TOC entry 3795 (class 0 OID 0)
-- Dependencies: 261
-- Name: position_to_workflow_mapping_position_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.position_to_workflow_mapping_position_id_seq', 1, false);


--
-- TOC entry 3796 (class 0 OID 0)
-- Dependencies: 262
-- Name: position_to_workflow_mapping_workflow_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.position_to_workflow_mapping_workflow_id_seq', 1, false);


--
-- TOC entry 3797 (class 0 OID 0)
-- Dependencies: 244
-- Name: rejected_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rejected_ticket_id_seq', 1, false);


--
-- TOC entry 3798 (class 0 OID 0)
-- Dependencies: 245
-- Name: rejected_ticket_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rejected_ticket_ticket_id_seq', 1, false);


--
-- TOC entry 3799 (class 0 OID 0)
-- Dependencies: 248
-- Name: section_status_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.section_status_ticket_id_seq', 3, true);


--
-- TOC entry 3800 (class 0 OID 0)
-- Dependencies: 229
-- Name: specified_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.specified_location_id_seq', 9, true);


--
-- TOC entry 3801 (class 0 OID 0)
-- Dependencies: 230
-- Name: specified_location_physical_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.specified_location_physical_location_id_seq', 1, false);


--
-- TOC entry 3802 (class 0 OID 0)
-- Dependencies: 217
-- Name: status_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.status_ticket_id_seq', 7, true);


--
-- TOC entry 3803 (class 0 OID 0)
-- Dependencies: 250
-- Name: status_ticket_section_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.status_ticket_section_id_seq', 1, false);


--
-- TOC entry 3804 (class 0 OID 0)
-- Dependencies: 233
-- Name: ticket_department_target_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ticket_department_target_id_seq', 1, false);


--
-- TOC entry 3805 (class 0 OID 0)
-- Dependencies: 232
-- Name: ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ticket_id_seq', 1, true);


--
-- TOC entry 3806 (class 0 OID 0)
-- Dependencies: 234
-- Name: ticket_physical_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ticket_physical_location_id_seq', 1, false);


--
-- TOC entry 3807 (class 0 OID 0)
-- Dependencies: 235
-- Name: ticket_specified_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ticket_specified_location_id_seq', 1, false);


--
-- TOC entry 3808 (class 0 OID 0)
-- Dependencies: 237
-- Name: track_status_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.track_status_ticket_id_seq', 1, true);


--
-- TOC entry 3809 (class 0 OID 0)
-- Dependencies: 239
-- Name: track_status_ticket_status_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.track_status_ticket_status_ticket_id_seq', 1, false);


--
-- TOC entry 3810 (class 0 OID 0)
-- Dependencies: 238
-- Name: track_status_ticket_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.track_status_ticket_ticket_id_seq', 1, false);


--
-- TOC entry 3811 (class 0 OID 0)
-- Dependencies: 254
-- Name: workflow_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.workflow_id_seq', 3, true);


--
-- TOC entry 3812 (class 0 OID 0)
-- Dependencies: 256
-- Name: workflow_step_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.workflow_step_id_seq', 6, true);


--
-- TOC entry 3813 (class 0 OID 0)
-- Dependencies: 258
-- Name: workflow_step_status_ticket_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.workflow_step_status_ticket_id_seq', 1, false);


--
-- TOC entry 3814 (class 0 OID 0)
-- Dependencies: 257
-- Name: workflow_step_workflow_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.workflow_step_workflow_id_seq', 1, false);


--
-- TOC entry 3451 (class 2606 OID 16431)
-- Name: area area_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.area
    ADD CONSTRAINT area_pkey PRIMARY KEY (id);


--
-- TOC entry 3453 (class 2606 OID 16433)
-- Name: area area_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.area
    ADD CONSTRAINT area_unique UNIQUE (name, department_id);


--
-- TOC entry 3447 (class 2606 OID 16415)
-- Name: department department_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.department
    ADD CONSTRAINT department_pkey PRIMARY KEY (id);


--
-- TOC entry 3449 (class 2606 OID 16417)
-- Name: department department_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.department
    ADD CONSTRAINT department_unique UNIQUE (name);


--
-- TOC entry 3455 (class 2606 OID 16452)
-- Name: employee employee_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT employee_pkey PRIMARY KEY (npk);


--
-- TOC entry 3471 (class 2606 OID 16572)
-- Name: job job_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT job_pkey PRIMARY KEY (id);


--
-- TOC entry 3473 (class 2606 OID 16574)
-- Name: job job_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT job_unique UNIQUE (ticket_id);


--
-- TOC entry 3501 (class 2606 OID 20537)
-- Name: permission permission_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permission
    ADD CONSTRAINT permission_pkey PRIMARY KEY (id);


--
-- TOC entry 3503 (class 2606 OID 20539)
-- Name: permission permission_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permission
    ADD CONSTRAINT permission_unique UNIQUE (name);


--
-- TOC entry 3457 (class 2606 OID 16474)
-- Name: physical_location physical_location_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.physical_location
    ADD CONSTRAINT physical_location_pkey PRIMARY KEY (id);


--
-- TOC entry 3459 (class 2606 OID 16476)
-- Name: physical_location physical_location_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.physical_location
    ADD CONSTRAINT physical_location_unique UNIQUE (name);


--
-- TOC entry 3505 (class 2606 OID 20552)
-- Name: position_permission position_permission_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.position_permission
    ADD CONSTRAINT position_permission_pkey PRIMARY KEY (position_id, permission_id);


--
-- TOC entry 3483 (class 2606 OID 18258)
-- Name: position position_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."position"
    ADD CONSTRAINT position_pkey PRIMARY KEY (id);


--
-- TOC entry 3497 (class 2606 OID 18334)
-- Name: position_to_workflow_mapping position_to_workflow_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.position_to_workflow_mapping
    ADD CONSTRAINT position_to_workflow_mapping_pkey PRIMARY KEY (id);


--
-- TOC entry 3499 (class 2606 OID 18347)
-- Name: position_to_workflow_mapping position_to_workflow_mapping_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.position_to_workflow_mapping
    ADD CONSTRAINT position_to_workflow_mapping_unique UNIQUE (position_id);


--
-- TOC entry 3485 (class 2606 OID 18260)
-- Name: position position_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."position"
    ADD CONSTRAINT position_unique UNIQUE (name);


--
-- TOC entry 3475 (class 2606 OID 16598)
-- Name: rejected_ticket rejected_ticket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rejected_ticket
    ADD CONSTRAINT rejected_ticket_pkey PRIMARY KEY (id);


--
-- TOC entry 3477 (class 2606 OID 18216)
-- Name: section_status_ticket section_status_ticket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.section_status_ticket
    ADD CONSTRAINT section_status_ticket_pkey PRIMARY KEY (id);


--
-- TOC entry 3479 (class 2606 OID 18220)
-- Name: section_status_ticket section_status_ticket_unique_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.section_status_ticket
    ADD CONSTRAINT section_status_ticket_unique_name UNIQUE (name);


--
-- TOC entry 3481 (class 2606 OID 18222)
-- Name: section_status_ticket section_status_ticket_unique_sequence; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.section_status_ticket
    ADD CONSTRAINT section_status_ticket_unique_sequence UNIQUE (sequence);


--
-- TOC entry 3461 (class 2606 OID 16490)
-- Name: specified_location specified_location_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.specified_location
    ADD CONSTRAINT specified_location_pkey PRIMARY KEY (id);


--
-- TOC entry 3463 (class 2606 OID 16492)
-- Name: specified_location specified_location_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.specified_location
    ADD CONSTRAINT specified_location_unique UNIQUE (physical_location_id, name);


--
-- TOC entry 3441 (class 2606 OID 16400)
-- Name: status_ticket status_ticket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.status_ticket
    ADD CONSTRAINT status_ticket_pkey PRIMARY KEY (id);


--
-- TOC entry 3443 (class 2606 OID 18224)
-- Name: status_ticket status_ticket_unique_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.status_ticket
    ADD CONSTRAINT status_ticket_unique_name UNIQUE (name);


--
-- TOC entry 3445 (class 2606 OID 18226)
-- Name: status_ticket status_ticket_unique_sequence; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.status_ticket
    ADD CONSTRAINT status_ticket_unique_sequence UNIQUE (sequence);


--
-- TOC entry 3465 (class 2606 OID 16514)
-- Name: ticket ticket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket
    ADD CONSTRAINT ticket_pkey PRIMARY KEY (id);


--
-- TOC entry 3467 (class 2606 OID 16548)
-- Name: track_status_ticket track_status_ticket_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket
    ADD CONSTRAINT track_status_ticket_id UNIQUE (ticket_id, status_ticket_id);


--
-- TOC entry 3469 (class 2606 OID 16546)
-- Name: track_status_ticket track_status_ticket_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket
    ADD CONSTRAINT track_status_ticket_pkey PRIMARY KEY (id);


--
-- TOC entry 3487 (class 2606 OID 18287)
-- Name: workflow workflow_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workflow
    ADD CONSTRAINT workflow_pkey PRIMARY KEY (id);


--
-- TOC entry 3491 (class 2606 OID 18304)
-- Name: workflow_step workflow_step_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workflow_step
    ADD CONSTRAINT workflow_step_pkey PRIMARY KEY (id);


--
-- TOC entry 3493 (class 2606 OID 18306)
-- Name: workflow_step workflow_step_unique_status_ticket_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workflow_step
    ADD CONSTRAINT workflow_step_unique_status_ticket_id UNIQUE (workflow_id, status_ticket_id);


--
-- TOC entry 3495 (class 2606 OID 18308)
-- Name: workflow_step workflow_step_unique_step_sequence; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workflow_step
    ADD CONSTRAINT workflow_step_unique_step_sequence UNIQUE (workflow_id, step_sequence);


--
-- TOC entry 3489 (class 2606 OID 18289)
-- Name: workflow workflow_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workflow
    ADD CONSTRAINT workflow_unique UNIQUE (name);


--
-- TOC entry 3530 (class 2620 OID 16610)
-- Name: area area_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER area_set_timestamp BEFORE UPDATE ON public.area FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3529 (class 2620 OID 16611)
-- Name: department department_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER department_set_timestamp BEFORE UPDATE ON public.department FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3531 (class 2620 OID 16612)
-- Name: employee employee; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER employee BEFORE UPDATE ON public.employee FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3535 (class 2620 OID 16613)
-- Name: job job_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER job_set_timestamp BEFORE UPDATE ON public.job FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3542 (class 2620 OID 20540)
-- Name: permission permission_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER permission_set_timestamp BEFORE UPDATE ON public.permission FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3532 (class 2620 OID 16614)
-- Name: physical_location physical_location_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER physical_location_set_timestamp BEFORE UPDATE ON public.physical_location FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3543 (class 2620 OID 20563)
-- Name: position_permission position_permission_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER position_permission_set_timestamp BEFORE UPDATE ON public.position_permission FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3538 (class 2620 OID 18274)
-- Name: position position_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER position_set_timestamp BEFORE UPDATE ON public."position" FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3541 (class 2620 OID 18345)
-- Name: position_to_workflow_mapping position_to_workflow_mapping_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER position_to_workflow_mapping_set_timestamp BEFORE UPDATE ON public.position_to_workflow_mapping FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3536 (class 2620 OID 16615)
-- Name: rejected_ticket rejected_ticket_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER rejected_ticket_set_timestamp BEFORE UPDATE ON public.rejected_ticket FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3537 (class 2620 OID 18242)
-- Name: section_status_ticket section_status_ticket_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER section_status_ticket_set_timestamp BEFORE UPDATE ON public.section_status_ticket FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3533 (class 2620 OID 16616)
-- Name: specified_location specified_location_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER specified_location_set_timestamp BEFORE UPDATE ON public.specified_location FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3528 (class 2620 OID 16617)
-- Name: status_ticket status_ticket_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER status_ticket_set_timestamp BEFORE UPDATE ON public.status_ticket FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3534 (class 2620 OID 16618)
-- Name: ticket ticket_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER ticket_set_timestamp BEFORE UPDATE ON public.ticket FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3539 (class 2620 OID 18290)
-- Name: workflow workflow_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER workflow_set_timestamp BEFORE UPDATE ON public.workflow FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3540 (class 2620 OID 18320)
-- Name: workflow_step workflow_step_set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER workflow_step_set_timestamp BEFORE UPDATE ON public.workflow_step FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- TOC entry 3508 (class 2606 OID 17889)
-- Name: employee area_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT area_id FOREIGN KEY (area_id) REFERENCES public.area(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3507 (class 2606 OID 17884)
-- Name: area department_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.area
    ADD CONSTRAINT department_id FOREIGN KEY (department_id) REFERENCES public.department(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3509 (class 2606 OID 17894)
-- Name: employee department_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT department_id FOREIGN KEY (department_id) REFERENCES public.department(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3512 (class 2606 OID 17924)
-- Name: ticket department_target_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket
    ADD CONSTRAINT department_target_id FOREIGN KEY (department_target_id) REFERENCES public.department(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3526 (class 2606 OID 20558)
-- Name: position_permission permission_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.position_permission
    ADD CONSTRAINT permission_id FOREIGN KEY (permission_id) REFERENCES public.permission(id) ON DELETE CASCADE;


--
-- TOC entry 3511 (class 2606 OID 17919)
-- Name: specified_location physical_location_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.specified_location
    ADD CONSTRAINT physical_location_id FOREIGN KEY (physical_location_id) REFERENCES public.physical_location(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3513 (class 2606 OID 17929)
-- Name: ticket physical_location_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket
    ADD CONSTRAINT physical_location_id FOREIGN KEY (physical_location_id) REFERENCES public.physical_location(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3518 (class 2606 OID 17899)
-- Name: job pic_job; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT pic_job FOREIGN KEY (pic_job) REFERENCES public.employee(npk) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3510 (class 2606 OID 18269)
-- Name: employee position_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee
    ADD CONSTRAINT position_id FOREIGN KEY (position_id) REFERENCES public."position"(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3524 (class 2606 OID 18335)
-- Name: position_to_workflow_mapping position_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.position_to_workflow_mapping
    ADD CONSTRAINT position_id FOREIGN KEY (position_id) REFERENCES public."position"(id) ON DELETE CASCADE;


--
-- TOC entry 3527 (class 2606 OID 20553)
-- Name: position_permission postion_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.position_permission
    ADD CONSTRAINT postion_id FOREIGN KEY (position_id) REFERENCES public."position"(id) ON DELETE CASCADE;


--
-- TOC entry 3520 (class 2606 OID 17909)
-- Name: rejected_ticket rejector; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rejected_ticket
    ADD CONSTRAINT rejector FOREIGN KEY (rejector) REFERENCES public.employee(npk) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3514 (class 2606 OID 17934)
-- Name: ticket requestor; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket
    ADD CONSTRAINT requestor FOREIGN KEY (requestor) REFERENCES public.employee(npk) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3506 (class 2606 OID 18243)
-- Name: status_ticket section_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.status_ticket
    ADD CONSTRAINT section_id FOREIGN KEY (section_id) REFERENCES public.section_status_ticket(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3515 (class 2606 OID 17939)
-- Name: ticket specified_location_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ticket
    ADD CONSTRAINT specified_location_id FOREIGN KEY (specified_location_id) REFERENCES public.specified_location(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3516 (class 2606 OID 17944)
-- Name: track_status_ticket status_ticket_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket
    ADD CONSTRAINT status_ticket_id FOREIGN KEY (status_ticket_id) REFERENCES public.status_ticket(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3522 (class 2606 OID 18314)
-- Name: workflow_step status_ticket_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workflow_step
    ADD CONSTRAINT status_ticket_id FOREIGN KEY (status_ticket_id) REFERENCES public.status_ticket(id) ON DELETE CASCADE;


--
-- TOC entry 3519 (class 2606 OID 17904)
-- Name: job ticket_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job
    ADD CONSTRAINT ticket_id FOREIGN KEY (ticket_id) REFERENCES public.ticket(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3521 (class 2606 OID 17914)
-- Name: rejected_ticket ticket_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rejected_ticket
    ADD CONSTRAINT ticket_id FOREIGN KEY (ticket_id) REFERENCES public.ticket(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3517 (class 2606 OID 17949)
-- Name: track_status_ticket ticket_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.track_status_ticket
    ADD CONSTRAINT ticket_id FOREIGN KEY (ticket_id) REFERENCES public.ticket(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3523 (class 2606 OID 18309)
-- Name: workflow_step workflow_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workflow_step
    ADD CONSTRAINT workflow_id FOREIGN KEY (workflow_id) REFERENCES public.workflow(id) ON DELETE CASCADE;


--
-- TOC entry 3525 (class 2606 OID 18340)
-- Name: position_to_workflow_mapping workflow_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.position_to_workflow_mapping
    ADD CONSTRAINT workflow_id FOREIGN KEY (workflow_id) REFERENCES public.workflow(id) ON DELETE CASCADE;


-- Completed on 2025-07-16 11:45:52

--
-- PostgreSQL database dump complete
--