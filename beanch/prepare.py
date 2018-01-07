#!/usr/bin/env python3
from lib.register_user import register_user
from lib.login_user import login_user
from lib.question import question
import config
import random

name = 'admin' + str(random.randrange(2 ** 32))
data = {
    'user_email': name + '@exmaple.com',
    'user_name': name,
    'user_password': name,
    'team_make': True,
    'team_name': name,
    'team_password': name,
}

register_user(config=config, data=data)
session = login_user(config=config, data=data).session()


q = question(config=config)

for _ in range(random.randint(30, 50)):
    q_add_data = {}
    q_add_data['name']  = 'problem %d' % random.randrange(2 ** 32)
    q_add_data['flag']  = 'HarekazeCTF{xxxxxx}'
    q_add_data['score'] = random.randint(1, 10) * 50
    q_add_data['genre'] = random.choice(['Pwn', 'Web', 'Rev', 'Crypto', 'Programming', 'Misc'])

    q.add(session=session, data=q_add_data)
