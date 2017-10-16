console.log( "==== sov2ex module: Filter ====" )

import TextField   from 'textfield';
import SelectField from 'selectfield';

const sort = [{
    value : "sumup",
    name  : "权重",
},{
    value : "created",
    name  : "发帖时间",
}],
order = [{
    value : "0",
    name  : "降序",
},{
    value : "1",
    name  : "升序",
}];

class Filter extends React.Component {

    state = {
        size_error : "",
        gte_error : "",
        lte_error : "",
        order_disable: true,
    };

    onSizeChange( event ) {
        const value = event.target.value.trim();
        if ( value == "" ) {
            this.setState({ size_error : "" });
            sessionStorage.removeItem( "size" );
        }
        else if ( !/^\d+$/.test( value ) || value < 1 || value > 50 ) {
            this.setState({
                size_error: "取值范围 1 ~ 50 的正整数"
            });
        } else {
            console.log( value )
            this.setState({ size_error : "" });
            sessionStorage.setItem( "size", value );
        }
    }

    onNodeChange( event ) {
        console.log( event.target.value.trim())
        event.target.value.trim() == "" ? sessionStorage.removeItem( "node" ) :
            sessionStorage.setItem( "node", event.target.value.trim() );
    }

    onSortChange( value, name ) {
        console.log( value, name )
        value == sort[0].value ? sessionStorage.removeItem( "sort" ) :
            sessionStorage.setItem( "sort", value );
         this.setState({
            order_disable: value == sort[0].value
        });
    }

    onOrderChange( value, name ) {
        console.log( value, name )
        value == order[0].value ? sessionStorage.removeItem( "order" ) :
            sessionStorage.setItem( "order", value );
    }

    getName( filter, value ) {
        if ( !value ) {
            return filter[0].name;
        } else {
            const result = filter.find( item => item.value == value );
            return result.name;
        }
    }

    getDay( value ) {
        if ( !value ) return "";
        else if ( !/\d+$/.test( value ) ) {
            return "";
        }
        else {
            const date   = new Date( parseInt( value )),
                  format = value => value = value < 10 ? "0" + value : value;
            return date.getFullYear() + "-" + format( date.getUTCMonth() + 1 ) + "-" + format( date.getUTCDate() );
        }
    }

    onDateChange( type, event ) {
        console.log( type, event.target.value )
        const value = event.target.value.trim(),
              error = `${type}_error`;
        if ( value == "" ) {
            this.setState({ [error] : "" });
            sessionStorage.removeItem( type );
        }
        else if ( /\w{4}-\w{2}-\w{2}/.test( value ) ) {
            const day = new Date( value ).getTime();
            if ( day ) {
                this.setState({ [error] : "" });
                sessionStorage.setItem( type, day );
            } else {
                this.setState({
                    [error]: "格式错误，如 2017-10-13"
                });
            }
        } else {
            this.setState({
                [error]: "格式错误，如 2017-10-13"
            });
        }
    }

    componentWillMount() {
        if ( location.search.startsWith( "?q=" ) ) {
            const query = window.location.search.replace( "?", "" ).split( "&" );
            query && query.length > 0 && query.forEach( item => {
                const [ key, value ] = item.split( "=" );
                key != "q" && sessionStorage.setItem( key, decodeURI( value ) );
            });
            console.log( sessionStorage )
        }
    }

    render() {
        return (
            <div className="filter">
                <TextField 
                    floatingtext="每页查询数量" placeholder="默认每页显示 10 条数据，取值范围在 1 ~ 50"
                    value={ sessionStorage.getItem( "size" ) }
                    errortext={ this.state.size_error }
                    onChange={ (e)=>this.onSizeChange(e) }
                />
                <TextField 
                    floatingtext="查询节点" placeholder="为空时，查询全部节点；支持节点名称与 节点 id"
                    value={ sessionStorage.getItem( "node" ) }
                    onChange={ (e)=>this.onNodeChange(e) }
                />
                <div className="horiz">
                    <TextField 
                        floatingtext="发帖的起始日期" placeholder="格式为 YYYY-MM-DD"
                        value={ this.getDay( sessionStorage.getItem( "gte" ) ) }
                        errortext={ this.state.gte_error }
                        onChange={ (evt)=>this.onDateChange( "gte", evt ) }
                    />
                    <TextField 
                        floatingtext="发帖的结束日期" placeholder="格式为 YYYY-MM-DD"
                        value={ this.getDay( sessionStorage.getItem( "lte" ) ) }
                        errortext={ this.state.lte_error }
                        onChange={ (evt)=>this.onDateChange( "lte", evt ) }
                    />
                </div>
                <div className="horiz">
                    <SelectField waves="md-waves-effect"
                        name={ this.getName( sort, sessionStorage.getItem( "sort" )) } items={ sort }
                        floatingtext="查询结果排序"
                        onChange={ (v,n)=>this.onSortChange(v,n) }
                    />
                    <SelectField waves="md-waves-effect"
                        disable={ !(sessionStorage.getItem( "sort" ) == sort[1].value) }
                        name={ this.getName( order, sessionStorage.getItem( "order" )) } items={ order }
                        floatingtext="发帖时间"
                        onChange={ (v,n)=>this.onOrderChange(v,n) }
                    />
                </div>
            </div>
        )
    }
}

function Render( target ) {
    ReactDOM.render( <Filter />, target );
}

export {
    Render,
}