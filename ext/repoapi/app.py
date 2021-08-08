from flask import Flask, session, redirect, url_for, escape, request, flash
from flask import render_template, abort, flash, g
from flask import jsonify
import json, requests
from functools import wraps
from operator import itemgetter, attrgetter
import time,os

from flask_cors import CORS, cross_origin

app = Flask(__name__)

cors = CORS(app)
app.config['CORS_HEADERS'] = 'Content-Type'

TOKEN = {'bearer_token': '', 'expired_time': 0}

def get_token():
    global TOKEN
    current_time = int(time.time())
    expired_time = TOKEN['expired_time']
    if current_time < expired_time:
        bearer_token = TOKEN['bearer_token']
        auz = {'Authorization': bearer_token}
        return auz
    base_url = 'https://bitbucket.org/site/oauth2/access_token'
    headers = {'Content-Type': 'application/x-www-form-urlencoded'}
    headers = ''
    params = {}
    data = {'grant_type': 'client_credentials'}
    client_id = os.environ['Key']
    secret = os.environ['Secret']
    req = requests.post(base_url, headers=headers, params=params, data=data, auth=(client_id, secret))
    req_dict = req.json()
    token = req_dict
    current_time = int(time.time())
    access_token = token.get('access_token')
    expires_in = token.get('expires_in')
    expired_time = current_time + expires_in - 600
    bearer_token = 'Bearer ' + access_token
    TOKEN = {'bearer_token': bearer_token, 'expired_time': expired_time}
    auz = {'Authorization': bearer_token}
    return auz

def is_json(myjson):
  try:
    json_object = json.loads(myjson)
  except ValueError as e:
    return False
  return True

@app.route('/')
def index():
    return 'ok'

@app.route('/console/v1/repoinfo')
@cross_origin()
def get_gitinfos():
    repoinfo = request.args.get('key')
    if repoinfo is None:
        ret = {
            'success': False,
            'results': "no repoinfo args"
        }
        return json.dumps(ret)
    repoinfo_detail = repoinfo.split("|")
    repoinfo_detail_len = len(repoinfo_detail)
    if repoinfo_detail_len < 3:
        ret = {
            'success': False,
            'results': "repoinfo args not enough"
        }
        return json.dumps(ret)
    if repoinfo_detail[0] != 'git' or repoinfo_detail[1] != 'bitbucket':
        ret = {
            'success': False,
            'results': "repoinfo args not correct"
        }
        return json.dumps(ret)
    reponame = repoinfo_detail[2]
    branches = get_branches(reponame)
    if isinstance(branches, (dict)):
        if branches.get('type') == 'error':
            return branches
    tags = get_tags(reponame)
    if isinstance(tags, (dict)):
        if tags.get('type') == 'error':
            return tags
    gitinfos = branches + tags
    ret = {
        'success': True,
        'results': gitinfos
    }
    return json.dumps(ret)

def get_branches(repo):
    base_url = 'https://api.bitbucket.org/2.0/repositories/'+repo+'/refs/branches'
    headers = {'Content-Type': 'application/json'}
    token = get_token()
    headers.update(token)
    params = {
        'pagelen': 100,
        'page': 1
    }
    data = '{}'
    branches = []
    while True:
        req = requests.get(base_url, headers=headers, params=params, data=data)

        branches_dict = req.json()
        if branches_dict.get('type') == 'error':
            return branches_dict
        branches = branches + branches_dict['values']
        if not branches_dict.get('next'):
            break
        params['page'] = params['page'] + 1
    branches = sorted(branches, key=lambda branch: branch['target']['date'].lower())
    new_branches =  []
    for branch in branches:
        new_branch = {}
        new_branch['name'] = branch['name'] + " " + branch['target']['hash'][0:7] + " " + branch['target']['date'] + " " + branch['target']['message']
        new_branch['value'] = 'branch/' + branch['name']
        new_branches.append(new_branch)
    return new_branches

def get_tags(repo):
    base_url = 'https://api.bitbucket.org/2.0/repositories/'+repo+'/refs/tags'
    headers = {'Content-Type': 'application/json'}
    token = get_token()
    headers.update(token)
    params = {
        'pagelen': 100
    }
    data = ''
    tags = []
    while True:
        req = requests.get(base_url, headers=headers, params=params, data=data)

        tags_dict = req.json()
        if tags_dict.get('type') == 'error':
            return tags_dict
        tags = tags + tags_dict['values']
        if not tags_dict.get('next'):
            break
        params['page'] = params['page'] + 1
    tags = sorted(tags, key=lambda tag: tag['target']['date'].lower())
    new_tags =  []
    for tag in tags:
        new_tag = {}
        new_tag['name'] = tag['name'] + " " + tag['target']['hash'][0:7] + " " + tag['target']['date'] + " " + tag['target']['message']
        new_tag['value'] = 'tag/' + tag['name']
        new_tags.append(new_tag)
    return new_tags
