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

#Generate People JSON objects
for i in range(0,1000):
    data={}
    data['first_name']=random.choice(fns)
    data['last_name']=random.choice(lns)

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
    data['first_name'] = name[0]
    data['last_name'] = name[1]
    data['title'] = fields[2]
    data['company_name'] = fields[4]
    people.append(data)


with open('people.json', 'w') as outfile:
     print "Saving People to file..."
     json.dump(people, outfile, indent=4, sort_keys=False)
#cols = ( 'object_id', 'field', 'transaction_type')
#
##dbname = os.environ['PG_DBNAME']
##user = os.environ['PG_USER']
#dbname = "njba"
#user = "manuel"
#connectTo = "dbname='" + dbname + "' user='" + user + "'"
#try:
#    con = pg.connect(connectTo)
#except:
#    print "Can't connect to db"
#    sys.exit(1)
#con.autocommit = True
#
#with open('obj_types.json') as data_file:
#    object_types = json.load(data_file)
#
#with open('objects.json') as data_file:
#    objects = json.load(data_file)
#
#with open('fields.json') as data_file:
#    fields = json.load(data_file)
#
#with open('values.json') as data_file:
#    values = json.load(data_file)
#
#def convert(name):
#    s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
#    return re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1).lower()
#
#obj_type_ids = []
#obj_ids = []
#
#cur = con.cursor()
#
#cur.execute("""TRUNCATE history CASCADE;""")
#cur.execute("""TRUNCATE transactions CASCADE;""")
#cur.execute("""TRUNCATE objects CASCADE;""")
#cur.execute("""TRUNCATE fields CASCADE;""")
#cur.execute("""TRUNCATE object_types CASCADE;""")
#
#def insert_rows(table_name, table_data):
#    columns = table_data[0].keys()
#    columns = ','.join([ convert(x) for x in columns ])
#    for obj in table_data:
#        values = "','".join(obj.values())
#        cur.execute("""INSERT into %s (%s) values ('%s');"""% (table_name, columns, values))
#
#def insert_transactions(values):
#    columns = ['object_id', 'field_id', 'value', 'transaction_type',  'created']
#    columns = ','.join([ convert(x) for x in columns ])
#    table_name = 'transactions'
#    for value in values:
#        strings = [value['object_id'], value['field']['id'], value['value'], 'update']
#        strings = "'" + "','".join(strings) + "'"
#        dates = str(value['date'])
#        dates = ", to_timestamp(" + dates + ")"
#        values = strings + dates
#        cur.execute("""INSERT into %s (%s) values (%s);"""% (table_name, columns, values))
#
#def insert_fields(table_data):
#    for obj in table_data:
#        cur.execute("""INSERT into fields (object_type_id, id, value_type,
#                title) values ('%s', '%s', %i, '%s');"""%
#                (obj['object_type_id'], obj['id'], obj['value_type'],
#                obj['title']))
#
#insert_rows('object_types', object_types)
#insert_rows('objects', objects)
#insert_fields(fields)
#insert_transactions(values)
