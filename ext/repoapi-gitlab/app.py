from flask import Flask, session, redirect, url_for, escape, request, flash
from flask import render_template, abort, flash, g
from flask import jsonify
import json, requests
from functools import wraps
from operator import itemgetter, attrgetter
import time,os

# from flask_cors import CORS, cross_origin

app = Flask(__name__)

# cors = CORS(app)
# app.config['CORS_HEADERS'] = 'Content-Type'

# TOKEN = {'bearer_token': '', 'expired_time': 0}
PRIVATE_TOKEN = 'glpat-RPLqsVwz-HhYjMMVe3BH'

def get_token():
    global PRIVATE_TOKEN
    auz = {'PRIVATE-TOKEN': PRIVATE_TOKEN}
    return auz

@app.route('/')
def index():
    return 'ok'

@app.route('/console/v1/repoinfo')
# @cross_origin()
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
    if repoinfo_detail[0] != 'git' or repoinfo_detail[1] != 'gitlab':
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
    base_url = 'https://git.yax.tech/api/v4/projects/'+repo+'/repository/branches'
    headers = {'Content-Type': 'application/json'}
    token = get_token()
    headers.update(token)
    params = {
        'per_page': 100,
        'page': 1
    }
    data = '{}'
    branches = []

    req = requests.get(base_url, headers=headers, params=params, data=data)

    branches = req.json()

    branches = sorted(branches, key=lambda branch: branch['commit']['committed_date'].lower())
    new_branches =  []
    for branch in branches:
        new_branch = {}
        new_branch['name'] = branch['name'] + " " + branch['commit']['short_id'] + " " + branch['commit']['committer_name'] + " " + branch['commit']['committed_date'] + " " + branch['commit']['title']
        new_branch['value'] = 'branch/' + branch['name']
        new_branches.append(new_branch)
    return new_branches

def get_tags(repo):
    base_url = 'https://git.yax.tech/api/v4/projects/'+repo+'/repository/tags'
    headers = {'Content-Type': 'application/json'}
    token = get_token()
    headers.update(token)
    params = {
        'per_page': 100,
        'page': 1
    }
    data = ''
    tags = []

    req = requests.get(base_url, headers=headers, params=params, data=data)

    tags = req.json()

    tags = sorted(tags, key=lambda tag: tag['commit']['committed_date'].lower())
    new_tags =  []
    for tag in tags:
        new_tag = {}
        new_tag['name'] = tag['name'] + " " + tag['commit']['short_id'] + " " + tag['commit']['committer_name'] + " " + tag['commit']['committed_date'] + " " + tag['commit']['title']
        new_tag['value'] = 'tag/' + tag['name']
        new_tags.append(new_tag)
    return new_tags
