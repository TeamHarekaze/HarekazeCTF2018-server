# File name           : login_user.py
# Author              : Hayato Doi
# Outline             : ログインの自動化を行うクラス
# license             : None
# Copyright (c) 2018, Hayato Doi

import requests
from bs4 import BeautifulSoup

class login_user:
    def __init__(self, config={},data={}):
        print('login_user(config={}, data={})'.format(config, data))
        self.target_url = config.target['url']

        self.user_email = data['user_email']
        self.user_password = data['user_password']

    def session(self):
        # make session
        session = requests.Session()

        # get token
        response = session.get( self.target_url + '/user/login' )
        if response.status_code != 200:
            raise Exception('fail')
        soup = BeautifulSoup( response.text.replace('</br>', ''), 'html.parser' )
        token = soup.body.find( 'form' ).find( attrs={ 'name' : 'csrf_token' } )['value']

        # post
        post_payload = {}
        post_payload['email']       = self.user_email
        post_payload['password']    = self.user_password
        post_payload['csrf_token']  = token

        response = session.post( self.target_url + '/user/login', data=post_payload)
        if response.status_code != 200:
            raise Exception('fail')

        return session