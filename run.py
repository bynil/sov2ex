import json
from flask import Flask, request, abort, Response, make_response, jsonify, render_template
from es_helper import es_search, es_clause_count
from datetime import datetime
from urllib3.exceptions import ReadTimeoutError
from elasticsearch import TransportError
from flask_moment import Moment

MAX_PAGING_DEPTH = 200
MAX_PAGING_SIZE = 50
MAX_KEYWORD_LENGTH = 100
MAX_CLAUSE_COUNT = 30
DEFAULT_SIZE = 10
DEFAULT_MAX_PAGE = int(MAX_PAGING_DEPTH / DEFAULT_SIZE)

moment = Moment()
app = Flask(__name__)
moment.init_app(app)


@app.route('/')
def index():
    keyword, page = parse_page_args(request)
    if not keyword:
        return render_template('index.html')

    es_size = DEFAULT_SIZE
    es_from = (page - 1) * DEFAULT_SIZE
    resp = web_search(keyword, es_from, es_size)

    total = resp['total']
    current_page = page
    max_page = min(DEFAULT_MAX_PAGE, int(total/DEFAULT_SIZE))
    pages = gen_pages(current_page, max_page)
    has_previous = current_page > 1
    has_next = current_page < max_page
    return render_template(
        'result.html', res=resp, pages=pages,
        current=current_page, q=keyword, enumerate=enumerate,
        has_previous=has_previous, has_next=has_next)


@app.route('/api/search', methods=['GET'])
def search_api():
    keyword, es_from, es_size = parse_api_args(request)
    if not keyword:
        abort(make_response(jsonify(message='Missing search keyword'), 400))
    resp = web_search(keyword, es_from, es_size)
    return Response(json.dumps(resp), mimetype='application/json')


def parse_page_args(req):
    """

    q: keyword
    page: page
    """
    try:
        keyword = req.args.get('q', None)
        page = int(req.args.get('page', 1))
        if page < 1:
            abort(make_response(jsonify(message='Wrong parameters'), 400))
        return keyword, page
    except ValueError:
        abort(make_response(jsonify(message='Wrong parameters'), 400))


def parse_api_args(req):
    """

    q: keyword
    from: from
    size: size
    """
    try:
        keyword = req.args.get('q', None)
        es_from = int(req.args.get('from', 0))
        es_size = int(req.args.get('size', DEFAULT_SIZE))
        if es_from < 0 or es_size < 0:
            abort(make_response(jsonify(message='Wrong parameters'), 400))
        return keyword, es_from, es_size
    except ValueError:
        abort(make_response(jsonify(message='Wrong parameters'), 400))


def web_search(keyword, es_from, es_size):
    if es_from + es_size > MAX_PAGING_DEPTH:
        abort(make_response(jsonify(message='Too deep paging parameters'), 400))
    if es_size > MAX_PAGING_SIZE:
        abort(make_response(jsonify(message='Too large size'), 400))
    try:
        if len(keyword) > MAX_KEYWORD_LENGTH or es_clause_count(keyword) > MAX_CLAUSE_COUNT:
            abort(make_response(jsonify(message='Too long keyword'), 400))
        result = es_search(keyword, es_from, es_size)
        hits = result['hits']
        resp = {'took': result['took'], 'timed_out': result['timed_out'],
                'total': hits['total'], 'hits': hits['hits']}
        return resp

    except ReadTimeoutError as e:
        msg = '{time} {url} {exception}'.format(time=datetime.now().isoformat(),
                                                url=request.url, exception=str(e))
        app.logger.error(msg)
        abort(make_response(jsonify(message='Read search result timeout'), 503))

    except TransportError as e:
        msg = '{time} {url} {exception}'.format(time=datetime.now().isoformat(),
                                                url=request.url, exception=str(e))
        app.logger.error(msg)
        abort(make_response(jsonify(message='Search engine error', detail=e.info), 503))

    except Exception as e:
        msg = '{time} {url} {exception}'.format(time=datetime.now().isoformat(),
                                                url=request.url, exception=str(e))
        app.logger.error(msg)
        abort(make_response(jsonify(message='Something went wrong'), 503))


def gen_pages(current, max_page, start_page=1, offset=3):

    """
        current = 4
        max_page = 8

        return [1 2 3 4 5 6 7]
    """
    interval = offset * 2 + 1  # 1....7
    if max_page <= interval + start_page - 1:
        return list(range(start_page, max_page + 1))

    if current - offset < start_page:
        return list(range(start_page, start_page + interval))

    if current + offset > max_page:
        return list(range(max_page - interval + 1, max_page + 1))

    return list(range(current - offset, current + offset + 1))


@app.template_filter('ctime')
def str2datetime(time_str):
    return datetime.strptime(time_str, '%Y-%m-%dT%H:%M:%S')

if __name__ == '__main__':
    app.run(host='127.0.0.1', debug=True)
