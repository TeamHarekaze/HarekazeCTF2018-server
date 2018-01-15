# File name           : register_user.py
# Author              : Hayato Doi
# Outline             : ユーザー登録の自動化を行うクラス
# license             : None
# Copyright (c) 2018, Hayato Doi

import requests
from bs4 import BeautifulSoup

class register_user:
    def __init__(self, config={},data={}):
        print('register_user(config={}, data={})'.format(config, data))
        self.target_url = config.target['url']

        self.user_email = data['user_email']
        self.user_name = data['user_name']
        self.user_password = data['user_password']

        self.team_join = data.get('team_join')
        self.team_make = data.get('team_make')

        self.team_name = data['team_name']
        self.team_password = data['team_password']

        # make session
        self.session = requests.Session()

        # get token
        response = self.session.get( self.target_url + '/user/register' )
        if response.status_code != 200:
            raise Exception('fail')
        soup = BeautifulSoup( response.text.replace('</br>', ''), 'html.parser' )
        token = soup.body.find( 'form' ).find( attrs={ 'name' : 'csrf_token' } )['value']

        # post
        post_payload = {}
        if self.team_join == True:
            post_payload['username']                        = self.user_name
            post_payload['email']                           = self.user_email
            post_payload['password']                        = self.user_password
            post_payload['password_confirmation']           = self.user_password
            post_payload['makejointeam']                    = 'join_team'
            post_payload['make_team_name']                  = ''
            post_payload['make_team_password']              = ''
            post_payload['make_team_password_confirmation'] = ''
            post_payload['join_team_name']                  = self.team_name
            post_payload['join_team_password']              = self.team_password
            post_payload['csrf_token']                      = token
        else:
            post_payload['username']                        = self.user_name
            post_payload['email']                           = self.user_email
            post_payload['password']                        = self.user_password
            post_payload['password_confirmation']           = self.user_password
            post_payload['makejointeam']                    = 'make_team'
            post_payload['make_team_name']                  = self.team_name
            post_payload['make_team_password']              = self.team_password
            post_payload['make_team_password_confirmation'] = self.team_password
            post_payload['join_team_name']                  = ''
            post_payload['join_team_password']              = ''
            post_payload['csrf_token']                      = token

        response = self.session.post( self.target_url + '/user/register', data=post_payload)
        if response.status_code != 200:
            print(response.text)
            raise Exception('fail')

    def session(self):
        return self.session