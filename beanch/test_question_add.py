from lib.login_user import login_user
from lib.question import question
import config

login_data = {
    'user_email':'user01@hoge.com',
    'user_password':'user01',
}
session = login_user(config=config, data=login_data).session()


q = question(config=config)

q_add_data = {}
q_add_data['name']  = 'hoge'
q_add_data['flag']  = 'HarekazeCTF{xxxxxx}'
q_add_data['score'] = 100
q_add_data['genre'] = 'web'

q.add(session=session, data=q_add_data)
