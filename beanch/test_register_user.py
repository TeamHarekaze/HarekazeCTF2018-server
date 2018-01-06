from lib.register_user import register_user
import config

register_data = {
    'user_email':'user01@hoge.com',
    'user_name':'user01',
    'user_password':'user01',
    'team_make':True,
    'team_name':'team01',
    'team_password':'team01',
}

register_user(config=config, data=register_data)