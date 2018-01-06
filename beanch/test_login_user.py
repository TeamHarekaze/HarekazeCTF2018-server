from lib.login_user import login_user
import config

login_data = {
    'user_email':'user01@hoge.com',
    'user_password':'user01',
}

session = login_user(config=config, data=login_data).session()
response = session.get( config.target['url'] + '/user/me' )
print(response.text)