console.log( "==== sov2ex Index ====" )

import './assets/css/style.css';
import './vender/notify/notify.css';

import Entry      from 'entry';
import Search     from 'search';
import Controlbar from 'controlbar';
import * as vers  from 'version';

import * as waves from 'waves';

vers.Init();
waves.Render({ root: "body" });
ReactDOM.render(
    location.search.startsWith( "?q=" ) ? <Search/> : <Entry/>,
    $( ".main" )[0]
);