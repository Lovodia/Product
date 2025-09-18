--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9
-- Dumped by pg_dump version 16.9

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, name) FROM stdin;
1	Товар 4
2	Товар 1
3	Товар 3
4	Товар 5
6	Электроника
7	Элек
8	жора
\.


--
-- Data for Name: goose_db_version; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.goose_db_version (id, version_id, is_applied, tstamp) FROM stdin;
1	0	t	2025-08-06 17:25:31.329104
2	20250806134613	t	2025-08-06 17:25:31.366189
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, name, price, category_id, created_at) FROM stdin;
1	MacBook Pro	1999.99	2	2025-08-06 18:08:41.136578
2	Pro	19.00	1	2025-08-06 18:09:08.232464
4	Шампунь для волос	799.99	2	2025-08-18 02:06:27.477135
5	Шампунь	79.00	3	2025-08-18 02:06:41.413378
6	нож	9.00	2	2025-08-18 02:07:28.526302
11	iPhone 14	956.00	4	2025-08-24 23:08:45.38008
12	iPhone 14	956.00	3	2025-08-24 23:08:53.137788
13	iPhone 112	6.00	2	2025-08-24 23:09:53.120037
\.


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 8, true);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.goose_db_version_id_seq', 2, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 13, true);


--
-- PostgreSQL database dump complete
--

