import json
from flask import Flask, request, abort, Response, make_response, jsonify, render_template
from es_helper import es_search, es_clause_count, es_time_order_search
from datetime import datetime
from urllib3.exceptions import ReadTimeoutError
from elasticsearch import TransportError
from flask_moment import Moment
from node_helper import find_node

MAX_PAGING_DEPTH = 1000
MAX_PAGING_SIZE = 50
MAX_KEYWORD_LENGTH = 100
MAX_CLAUSE_COUNT = 30
DEFAULT_SIZE = 10
DEFAULT_MAX_PAGE = int(MAX_PAGING_DEPTH / DEFAULT_SIZE)

moment = Moment()
app = Flask(__name__)
moment.init_app(app)

SUMUP_SORT = 'sumup'
CREATED_SORT = 'created'

sort_choice = [SUMUP_SORT, CREATED_SORT]
order_choice = [0, 1]


class QueryParams(object):
    def __init__(self, keyword, es_from, es_size, sort, order, gte, lte, node_id, operator):
        self.keyword = keyword
        self.es_from = es_from
        self.es_size = es_size
        self.sort = sort
        self.order = order
        self.gte = gte
        self.lte = lte
        self.node_id = node_id
        self.page = 1
        self.operator = operator


@app.route('/api/search', methods=['GET'])
def search_api():
    params = parse_api_args(request)
    if not params.keyword:
        abort(make_response(jsonify(message='Missing search keyword'), 400))
    resp = web_search(params)
    return Response(json.dumps(resp), mimetype='application/json')


def parse_api_args(req) -> QueryParams:
    """

    q: keyword
    from: from
    size: size
    sort: sort by
    order: order by asc or desc (not used when sort is `sumup`)
    gte: created time great than or equal to timestamp (epoch_second)
    lte: created time less than or equal to timestamp (epoch_second)
    node: limitative node
    operator: `and` or `or`, default is `or`
    """
    try:
        keyword = req.args.get('q', None)
        es_from = int(req.args.get('from', 0))
        es_size = int(req.args.get('size', DEFAULT_SIZE))
        if es_from < 0 or es_size < 0:
            abort(make_response(jsonify(message='Wrong parameters'), 400))

        sort = req.args.get('sort', sort_choice[0])
        if sort not in sort_choice:
            abort(make_response(jsonify(message='Wrong parameters'), 400))

        order = int(req.args.get('order', order_choice[0]))
        if order not in order_choice:
            abort(make_response(jsonify(message='Wrong parameters'), 400))

        gte = req.args.get('gte', None)
        lte = req.args.get('lte', None)
        if gte is not None:
            gte = int(gte)
        if lte is not None:
            lte = int(lte)

        node_name = req.args.get('node', None)
        node = find_node(node_name)
        if node:
            node_id = node['id']
        else:
            node_id = None

        operator = req.args.get('operator', 'or')
        if operator not in ['or', 'and']:
            operator = 'or'

        return QueryParams(keyword=keyword, es_from=es_from, es_size=es_size,
                           sort=sort, order=order, gte=gte, lte=lte, node_id=node_id, operator=operator)
    except ValueError:
        abort(make_response(jsonify(message='Wrong parameters'), 400))


def web_search(params: QueryParams):
    keyword, es_from, es_size, \
    sort, order, gte, lte, node_id, operator = params.keyword, params.es_from, params.es_size, params.sort, \
                            params.order, params.gte, params.lte, params.node_id, params.operator

    if es_from + es_size > MAX_PAGING_DEPTH:
        abort(make_response(jsonify(message='Too deep paging parameters'), 400))
    if es_size > MAX_PAGING_SIZE:
        abort(make_response(jsonify(message='Too large size'), 400))
    try:
        if len(keyword) > MAX_KEYWORD_LENGTH or es_clause_count(keyword) > MAX_CLAUSE_COUNT:
            abort(make_response(jsonify(message='Too long keyword'), 400))

        if sort == SUMUP_SORT:
            result = es_search(keyword, es_from, es_size, gte, lte, node_id, operator)
        elif sort == CREATED_SORT:
            result = es_time_order_search(keyword, es_from, es_size, order, gte, lte, node_id, operator)
        else:
            result = es_search(keyword, es_from, es_size, gte, lte, node_id, operator)

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


if __name__ == '__main__':
    app.run(host='127.0.0.1', debug=True)
