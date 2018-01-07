from lib.register_user import register_user
import config
import time
from statistics import mean

time_list = []
fail_count = 0
for i in range(10):
    start_time = time.time()
    _fail_count_tmp = fail_count
    register_team_data = {
        'user_email':'team' + str(i) + 'user' + str(0) +'@hoge.com',
        'user_name':'team' + str(i) + 'user' + str(0) ,
        'user_password':'team' + str(i) + 'user' + str(0) ,
        'team_make':True,
        'team_name':'team' + str(i) ,
        'team_password':'team' + str(i),
    }
    print( 'user' + str(0) +'@hoge.com')
    try:
        register_user(config=config, data=register_team_data)
    except Exception:
        print('fail : make team')
        fail_count = fail_count + 1
        break
    for j in range(1, 10):
        register_user_data = {
            'user_email':'team' + str(i) + 'user' + str(j) +'@hoge.com',
            'user_name':'team' + str(i) + 'user' + str(j) ,
            'user_password':'team' + str(i) + 'user' + str(j) ,
            'team_join':True,
            'team_name':'team' + str(i) ,
            'team_password':'team' + str(i),
        }
        print('user' + str(j) +'@hoge.com')
        try:
            register_user(config=config, data=register_user_data)
        except Exception:
            print('fail : make user')
            fail_count = fail_count + 1
            break
    if _fail_count_tmp == fail_count:
        time_list.append( time.time() - start_time)

print('count :' + str(len(time_list) + fail_count))
print('time max :' + str(max(time_list)))
print('time min :' + str(min(time_list)))
print('time ave :' + str(sum(time_list) / len(time_list)))
print('fail :' + str(fail_count / ( len(time_list) + fail_count )) + '%')
