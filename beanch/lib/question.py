# File name           : question.py
# Author              : Hayato Doi
# Outline             : 問題関係の自動化を行うクラス
# license             : None
# Copyright (c) 2018, Hayato Doi

import requests
from bs4 import BeautifulSoup

class question:
    def __init__(self, config={}):
        self.target_url = config.target['url']
        self.target_admin_url = config.target['admin_url']

    # _re_token : HTMLを受け取り、解析し、tokenを返す
    def _re_token(self, html):
        soup = BeautifulSoup( html.replace('</br>', ''), 'html.parser' )
        token = soup.body.find( 'form' ).find( attrs={ 'name' : 'csrf_token' } )['value']
        return token

    # add : 問題を追加する
    def add(self, session, data={}):
        print('question.add(data={})'.format(data))
        # get token
        response = session.get( self.target_admin_url + '/question/add' )
        if response.status_code != 200:
            raise Exception('fail')
        token = self._re_token(response.text)

        # post
        post_payload = {}
        post_payload['name']               = data['name']
        post_payload['flag']               = data['flag']
        post_payload['score']              = data['score']
        post_payload['genre']              = data['genre']
        post_payload['publish_now']        = 'on'
        post_payload['publish_start_time'] = ''
        post_payload['sentence']           = 'Flag is ' + data['flag']
        post_payload['csrf_token']         = token

        response = session.post( self.target_admin_url + '/question/add', data=post_payload)
        if response.status_code != 200:
            raise Exception('fail')

    # submit : 問題に回答する 戻り値:正解(True),不正解(False)
    def submit(self, session, q_id=1, flag=''):
        print('question.submit(session={}, q_id={}, flag={})'.format(session, q_id, flag))
        # get token
        response = session.get( self.target_url + '/answer/' + str(q_id) )
        if response.status_code != 200:
            raise Exception('fail')
        token = self._re_token(response.text)

        post_payload = {}
        post_payload['flag']       = flag
        post_payload['csrf_token'] = token

        response = session.post( self.target_url + '/answer/' + str(q_id), data=post_payload)
        if response.status_code != 200:
            raise Exception('fail')
        soup = BeautifulSoup( response.text.replace('</br>', ''), 'html.parser' )
        msg = soup.body.find( 'form' ).find( attrs={ 'role' : 'alert' } ).string
        if msg == 'Correct answer':
            return True
        else:
            return False

    # list : 問題一覧を習得する 戻り値:id list
    def list(self):
        print('question.list()')
        q_id_list = []
        session = requests.Session()
        
        response = session.get( self.target_url + '/question' )
        if response.status_code != 200:
            raise Exception('fail')
        soup = BeautifulSoup( response.text.replace('</br>', ''), 'html.parser' )
        for i in soup.find('tbody').findAll('tr'):
            url = i.find('a')['href']
            url_sp = url.split("/")
            q_id_list.append(url_sp[len(url_sp) - 1])

        return q_id_list