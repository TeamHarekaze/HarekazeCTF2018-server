#!/usr/bin/env python3
from lib.login_user import login_user
from lib.register_user import register_user
from lib.question import question
import config
import random

import argparse
parser = argparse.ArgumentParser()
parser.add_argument('n', type=int)
args = parser.parse_args()

def make_data(team, user, team_make=False):
    return {
        'user_email': user + '@exmaple.com',
        'user_name': user,
        'user_password': user,
        'team_make': team_make,
        'team_join': not team_make,
        'team_name': team,
        'team_password': team,
    }

for _ in range(args.n):
    team = 'team' + str(random.randrange(2 ** 64))
    sessions = []
    for i in range(random.randint(1, 6)):
        user = 'member' + str(random.randrange(2 ** 64))
        data = make_data(team, user, team_make=not i)
        register_user(config=config, data=data)
        session = login_user(config=config, data=data).session()
        sessions += [ session ]

    q = question(config=config)
    q_ids = list(range(1, 40))  # q.list()  # workaround for issue/109
    random.shuffle(q_ids)
    q_ids = q_ids[: random.randint(1, len(q_ids))]
    for q_id in q_ids:
        session = random.choice(sessions)
        q.submit(session=session, q_id=q_id, flag='HarekazeCTF{xxxxxx}')
        if random.random() < 0.3:
            # q.list()  # workaround
            print('session.get(/question)')  # workaround
            session.get( config.target['url'] + '/question' )  # workaround
            pass
        if random.random() < 0.3:
            print('session.get(/ranking)')
            session.get( config.target['url'] + '/ranking' )
        if random.random() < 0.3:
            print('session.get(/user/me)')
            session.get( config.target['url'] + '/user/me' )
