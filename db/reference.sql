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
