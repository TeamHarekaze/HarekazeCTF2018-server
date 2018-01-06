from lib.login_user import login_user
from lib.question import question
import config

login_data = {
    'user_email':'user01@hoge.com',
    'user_password':'user01',
}
session = login_user(config=config, data=login_data).session()


q = question(config=config)

if q.submit(session=session,  q_id=1, flag='HarekazeCTF{Question1}'):
    print('正解しました。')
else:
    print('不正解でした。')