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

for _ in range(args.n):
    name = 'stress' + str(random.randrange(2 ** 64))
    data = {
        'user_email': name + '@exmaple.com',
        'user_name': name,
        'user_password': name,
        'team_make': True,
        'team_name': name,
        'team_password': name,
    }

    # TODO: チームメンバーを複数にする

    register_user(config=config, data=data)
    session = login_user(config=config, data=data).session()
    response = session.get( config.target['url'] + '/user/me' )
    q = question(config=config)
    q_ids = q.list()
    random.shuffle(q_ids)
    q_ids = q_ids[: random.randint(3, 30)]
    for q_id in q_ids:
        q.submit(session=session, q_id=q_id, flag='HarekazeCTF{xxxxxx}')
        if random.random() < 0.3:
            q.list()
        if random.random() < 0.3:
            session.get( config.target['url'] + '/ranking' )
