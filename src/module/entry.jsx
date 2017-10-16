console.log( "==== sov2ex module: Entry ====" )

import TextField   from 'textfield';

import * as filter from 'filter';

export default class Entry extends React.Component{

    onKeyDown( event ) {
        event.keyCode == 13 &&
            this.search( event.target.value );
    }

    onClick() {
        this.search( this.refs.search.refs.target.value );
    }

    arrowOnClick() {
        filter.Render( $( ".filtergp" )[0] );
        $( ".filtergp" ).toggleClass( "filtergp-top" );
    }

    search( value ) {
        if ( value.trim() != "" ) {
            let url = window.location.origin + window.location.pathname + `?q=${value}`;
            Object.keys( sessionStorage ).forEach( key => url += `&${key}=${sessionStorage[key]}`);
            sessionStorage.clear();
            console.log( sessionStorage, url )
            window.location.href = url;
        } else {
            new Notify().Render( "不能为空，请输入正确的值。" );
        }
    }

    componentWillMount() {
        sessionStorage.clear();
    }

    render() {
        return (
            <div className="entry">
                <div className="logo">
                    <img src="./assets/images/logo@2x.png"></img>
                </div>
                <div className="searchbar">
                    <div className="search">
                        <TextField 
                            ref="search"
                            placeholder="请输入查询的关键字" 
                            onKeyDown={ (e)=> this.onKeyDown(e) }
                        />
                        <span className="bar" onClick={ ()=> this.onClick() }></span>
                        <span className="arrow" onClick={ ()=> this.arrowOnClick() }></span>
                        <div className="filtergp"></div>
                    </div>
                </div>
                <div className="desc">
                    一个便捷的 v2ex 站内搜索引擎，<a href="https://github.com/Bynil/sov2ex" target="_blank">了解更多</a>  。
                </div>
            </div>
        )
    }
}