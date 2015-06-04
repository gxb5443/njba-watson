DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS people(
  id uuid PRIMARY KEY not null default uuid_generate_v4(),
  name text not null,
  address1 text,
  address2 text,
  zip varchar(5),
  state varchar(2),
  country varchar(40),
  company_id uuid[],
  home_phone varchar(10),
  work_phone varchar(10),
  cell_phone varchar(10),
  email_address text,
  member_id uuid,
  created timestamp without time zone default(now() at time zone 'utc')
);

CREATE TABLE IF NOT EXISTS companies (
  id uuid PRIMARY KEY not null default uuid_generate_v4(),
  name text not null,
  address1 text,
  address2 text,
  zip varchar(5),
  state varchar(2),
  country varchar(40),
  parent uuid references companies(id),
  poc uuid references people(id),
  created timestamp without time zone default(now() at time zone 'utc')
);

CREATE SEQUENCE member_id_seq;
CREATE TABLE IF NOT EXISTS memberships(
  id smallint PRIMARY KEY NOT NULL DEFAULT nextval('member_id_seq'),
  status varchar(15) NOT NULL DEFAULT 'active',
  created timestamp without time zone default(now() at time zone 'utc')
);
ALTER SEQUENCE member_id_seq OWNED BY memberships.id;
SELECT setval('member_id_seq', 10000);

CREATE TABLE IF NOT EXISTS membership_holdings (
  member_id smallint not null references memberships(id),
  company_id uuid not null references companies(id),
  duration daterange NOT NULL DEFAULT '[today, infinity)'::daterange,
  created timestamp without time zone default(now() at time zone 'utc')
);

CREATE TABLE IF NOT EXISTS groups (
  id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  name text NOT NULL,
  description text,
  duration daterange NOT NULL DEFAULT '[today, infinity)'::daterange,
  created timestamp without time zone default(now() at time zone 'utc')
);

CREATE TABLE IF NOT EXISTS group_membership (
  group_id uuid not null references groups(id),
  user_id uuid not null references people(id),
  role text,
  duration daterange NOT NULL DEFAULT '[today, infinity)'::daterange,
  created timestamp without time zone default(now() at time zone 'utc')
);

INSERT INTO PEOPLE (name) VALUES ('Gian Biondi');
INSERT INTO PEOPLE (name) VALUES ('Mike Biondi');
INSERT INTO COMPANIES (name, poc) VALUES ('Umbrella Co', (SELECT id from people where name='Gian Biondi'));
INSERT INTO COMPANIES (name, parent, poc) VALUES ('Child Co', (SELECT id from companies where name='Umbrella Co'), (SELECT id from people where name='Mike Biondi'));
INSERT INTO COMPANIES (name, parent, poc) VALUES ('Sibling Co', (SELECT id from companies where name='Umbrella Co'), (SELECT id from people where name='Gian Biondi'));
INSERT INTO memberships default values;
INSERT INTO membership_holdings (member_id, company_id) VALUES ((SELECT id from memberships limit 1), (SELECT id FROM companies WHERE name='Umbrella Co'));

INSERT INTO groups (name) VALUES ('BoD');
INSERT INTO group_membership (group_id, user_id, role) VALUES ((SELECT id from groups where name='BoD' LIMIT 1),(SELECT id from people where name = 'Gian Biondi'), 'President');
