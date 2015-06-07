import random
import re
import json
import psycopg2 as pg
import sys
import uuid
import os

random.seed()

fns=[
    'Gian',
    'Grant',
    'Lisa',
    'Diane',
    'Pauline',
    'Carol Ann',
    'Jeff',
    'Cindy'
    ]
lns=[
    'Smith',
    'Destine',
    'DePastene',
    'Nicolo',
    'Pocino',
    'Spicer',
    'Lucking',
    'Spalding',
    'Dandy'
    ]
peopleByCo = {}
coById = {}

fp = open('dcalists/NJ_BLDR_PRIN_20141113.txt', 'r')
next(fp)
people = []
for line in fp.xreadlines():
    data={}
    fields = line.split('|')
    fname = ''.join(i for i in fields[1] if not i.isdigit()).strip()
    name = fname.split(' ', 1)
    if len(name)<2:
        continue
    data['first_name'] = name[0].lower().title()
    data['last_name'] = name[1].lower().title()
    data['title'] = fields[2]
    data['company_name'] = fields[4]
    data['id'] = str(uuid.uuid4())
    people.append(data)
    peopleByCo[fields[4]] = data['id']

fc = open('dcalists/NJ_BLDR_20141113.txt', 'r')
next(fc)
companies= []
for line in fc.xreadlines():
    data={}
    fields = line.split('|')
    data['company_name'] = fields[1]
    data['address1'] = fields[2]
    data['address2'] = fields[3]
    data['city'] = fields[4]
    data['state'] = fields[5]
    data['zip'] = fields[6]
    data['id'] = str(uuid.uuid4())
    companies.append(data)
    coById[fields[1]] = data['id']

with open('companies.json', 'w') as outfile:
     print "Saving Companies to file..."
     json.dump(companies, outfile, indent=4, sort_keys=False)

with open('people.json', 'w') as outfile:
     print "Saving People to file..."
     json.dump(people, outfile, indent=4, sort_keys=False)

try:
    con = pg.connect("dbname='njba' user='postgres'")
except:
    print "Can't connect to db"
    sys.exit(1)
con.autocommit=True
cur = con.cursor()

print "Loading companies to db..."
for c in companies:
    cur.execute("""INSERT into companies (id, name, address1, address2, zip, state, city) VALUES (%s, %s, %s, %s, %s, %s, %s)""", (c['id'], c['company_name'], c['address1'], c['address2'], c['zip'], c['state'], c['city']))

print "Loading people to db..."
for p in people:
    #if p['company_name'] in coById:
    #    p['company_id'] = coById[p['company_name']]
    cur.execute("""INSERT into people (id, first_name, last_name, title) values (%s, %s, %s, %s);""", (p["id"], p['first_name'], p['last_name'], p['title']))
    if p['company_name'] in coById:
        cur.execute("""INSERT into pco_relationships (people_id, company_id) VALUES (%s, %s)""", (p['id'], coById[p['company_name']]))
