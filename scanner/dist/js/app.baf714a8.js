(function(t){function e(e){for(var n,s,i=e[0],c=e[1],u=e[2],d=0,f=[];d<i.length;d++)s=i[d],a[s]&&f.push(a[s][0]),a[s]=0;for(n in c)Object.prototype.hasOwnProperty.call(c,n)&&(t[n]=c[n]);l&&l(e);while(f.length)f.shift()();return o.push.apply(o,u||[]),r()}function r(){for(var t,e=0;e<o.length;e++){for(var r=o[e],n=!0,i=1;i<r.length;i++){var c=r[i];0!==a[c]&&(n=!1)}n&&(o.splice(e--,1),t=s(s.s=r[0]))}return t}var n={},a={app:0},o=[];function s(e){if(n[e])return n[e].exports;var r=n[e]={i:e,l:!1,exports:{}};return t[e].call(r.exports,r,r.exports,s),r.l=!0,r.exports}s.m=t,s.c=n,s.d=function(t,e,r){s.o(t,e)||Object.defineProperty(t,e,{enumerable:!0,get:r})},s.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},s.t=function(t,e){if(1&e&&(t=s(t)),8&e)return t;if(4&e&&"object"===typeof t&&t&&t.__esModule)return t;var r=Object.create(null);if(s.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:t}),2&e&&"string"!=typeof t)for(var n in t)s.d(r,n,function(e){return t[e]}.bind(null,n));return r},s.n=function(t){var e=t&&t.__esModule?function(){return t["default"]}:function(){return t};return s.d(e,"a",e),e},s.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},s.p="";var i=window["webpackJsonp"]=window["webpackJsonp"]||[],c=i.push.bind(i);i.push=e,i=i.slice();for(var u=0;u<i.length;u++)e(i[u]);var l=c;o.push([0,"chunk-vendors"]),r()})({0:function(t,e,r){t.exports=r("56d7")},"07bf":function(t,e,r){"use strict";var n=r("245a"),a=r.n(n);a.a},"17b1":function(t,e,r){"use strict";var n=r("1f24"),a=r.n(n);a.a},"1f24":function(t,e,r){},"245a":function(t,e,r){},"245e":function(t,e,r){"use strict";var n=r("5999"),a=r.n(n);a.a},2906:function(t,e,r){"use strict";var n=r("f13f"),a=r.n(n);a.a},"3bfe":function(t,e,r){},"474f":function(t,e,r){"use strict";var n=r("9c5c"),a=r.n(n);a.a},"4def":function(t,e,r){},"537b":function(t,e,r){"use strict";var n=r("3bfe"),a=r.n(n);a.a},"56d7":function(t,e,r){"use strict";r.r(e);r("744f"),r("6c7b"),r("7514"),r("20d6"),r("1c4c"),r("6762"),r("cadf"),r("e804"),r("55dd"),r("d04f"),r("c8ce"),r("217b"),r("7f7f"),r("f400"),r("7f25"),r("536b"),r("d9ab"),r("f9ab"),r("32d7"),r("25c9"),r("9f3c"),r("042e"),r("c7c6"),r("f4ff"),r("049f"),r("7872"),r("a69f"),r("0b21"),r("6c1a"),r("c7c62"),r("84b4"),r("c5f6"),r("2e37"),r("fca0"),r("7cdf"),r("ee1d"),r("b1b1"),r("87f3"),r("9278"),r("5df2"),r("04ff"),r("f751"),r("4504"),r("fee7"),r("ffc1"),r("0d6d"),r("9986"),r("8e6e"),r("25db"),r("e4f7"),r("b9a1"),r("64d5"),r("9aea"),r("db97"),r("66c8"),r("57f0"),r("165b"),r("456d"),r("cf6a"),r("fd24"),r("8615"),r("551c"),r("097d"),r("df1b"),r("2397"),r("88ca"),r("ba16"),r("d185"),r("ebde"),r("2d34"),r("f6b3"),r("2251"),r("c698"),r("a19f"),r("9253"),r("9275"),r("3b2b"),r("3846"),r("4917"),r("a481"),r("28a5"),r("386d"),r("6b54"),r("4f7f"),r("8a81"),r("ac4d"),r("8449"),r("9c86"),r("fa83"),r("48c0"),r("a032"),r("aef6"),r("d263"),r("6c37"),r("9ec8"),r("5695"),r("2fdb"),r("d0b0"),r("5df3"),r("b54a"),r("f576"),r("ed50"),r("788d"),r("14b9"),r("f386"),r("f559"),r("1448"),r("673e"),r("242a"),r("c66f"),r("b05c"),r("34ef"),r("6aa2"),r("15ac"),r("af56"),r("b6e4"),r("9c29"),r("63d9"),r("4dda"),r("10ad"),r("c02b"),r("4795"),r("130f"),r("ac6a"),r("96cf"),r("0cdd");var n=r("2b0e"),a=r("bc3a"),o=r.n(a),s={},i=o.a.create(s);i.interceptors.request.use(function(t){return t},function(t){return Promise.reject(t)}),i.interceptors.response.use(function(t){return t},function(t){return Promise.reject(t)}),Plugin.install=function(t){t.axios=i,window.axios=i,Object.defineProperties(t.prototype,{axios:{get:function(){return i}},$axios:{get:function(){return i}}})},n["default"].use(Plugin);Plugin;var c=r("9f7b"),u=r.n(c);r("ab8b"),r("2dd8");n["default"].use(u.a);var l=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{attrs:{id:"main"}},[t.fatal?r("FatalError",{attrs:{error:t.fatal}}):t.auth?t.event?r("TicketScanner",{attrs:{auth:t.auth,event:t.event,freeEntry:t.event.freeEntries}}):r("EventChooser",{attrs:{auth:t.auth},on:{event:function(e){t.event=e},fatal:function(e){t.fatal=e}}}):r("LoginPage",{on:{auth:function(e){t.auth=e},fatal:function(e){t.fatal=e}}})],1)},d=[],f=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{attrs:{id:"events"}},[r("LogoWide"),t.events.length?r("div",{attrs:{id:"events-head"}},[t._v("Which event are you taking tickets for?")]):t._e(),t._l(t.events,function(e){return r("div",{key:e.id,staticClass:"event",on:{click:function(r){return t.$emit("event",e)}}},[t._v(t._s(e.start.substr(0,10))+" "+t._s(e.name))])})],2)},h=[],p=r("a34a"),v=r.n(p),m=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{attrs:{id:"logo-wide"}},[t._v("Schola Cantorum Ticket Scanner")])},b=[],A=(r("c3d6"),r("2877")),x={},g=Object(A["a"])(x,m,b,!1,null,null,null),w=g.exports;function y(t,e,r,n,a,o,s){try{var i=t[o](s),c=i.value}catch(u){return void r(u)}i.done?e(c):Promise.resolve(c).then(n,a)}function k(t){return function(){var e=this,r=arguments;return new Promise(function(n,a){var o=t.apply(e,r);function s(t){y(o,n,a,s,i,"next",t)}function i(t){y(o,n,a,s,i,"throw",t)}s(void 0)})}}var E={components:{LogoWide:w},props:{auth:String},data:function(){return{events:[]}},mounted:function(){var t=k(v.a.mark(function t(){var e,r=this;return v.a.wrap(function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,this.$axios.get("/api/event?future=1&freeEntries=1",{headers:{auth:this.auth}}).catch(function(t){return console.log(t),r.$emit("fatal","Server error"),null});case 2:if(e=t.sent,e){t.next=5;break}return t.abrupt("return");case 5:if(200===e.status){t.next=9;break}return console.log(e.statusText),this.$emit("fatal","Server error"),t.abrupt("return");case 9:this.events=e.data;case 10:case"end":return t.stop()}},t,this)}));function e(){return t.apply(this,arguments)}return e}()},S=E,P=(r("a4ff"),Object(A["a"])(S,f,h,!1,null,null,null)),_=P.exports,O=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{attrs:{id:"fatal-top"}},[r("LogoWide"),r("div",{attrs:{id:"fatal"}},[r("div",{domProps:{textContent:t._s(t.error)}}),r("b-btn",{attrs:{variant:"primary"},on:{click:t.onReload}},[t._v("Reload")])],1)],1)},C=[],T={components:{LogoWide:w},props:{error:String},methods:{onReload:function(){location.reload()}}},j=T,L=(r("58b5"),Object(A["a"])(j,O,C,!1,null,null,null)),R=L.exports,V=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{attrs:{id:"login"}},[r("LogoWide"),r("div",{attrs:{id:"login-head"}},[t._v("Please log in.")]),r("form",{attrs:{id:"login-form"},on:{submit:function(e){return e.preventDefault(),t.onSubmit(e)}}},[r("b-form-group",{attrs:{inline:"",label:"Username","label-cols":"4","label-for":"login-username"}},[r("b-form-input",{attrs:{id:"login-username",autofocus:""},model:{value:t.username,callback:function(e){t.username=e},expression:"username"}})],1),r("b-form-group",{attrs:{inline:"",label:"Password","label-cols":"4","label-for":"login-password"}},[r("b-form-input",{attrs:{id:"login-password",type:"password"},model:{value:t.password,callback:function(e){t.password=e},expression:"password"}})],1),r("b-btn",{attrs:{id:"login-btn",type:"submit",variant:"primary"}},[t._v("Login")])],1),r("div",{directives:[{name:"show",rawName:"v-show",value:t.error,expression:"error"}],attrs:{id:"login-error"},domProps:{textContent:t._s(t.error)}})],1)},D=[];function F(t,e,r,n,a,o,s){try{var i=t[o](s),c=i.value}catch(u){return void r(u)}i.done?e(c):Promise.resolve(c).then(n,a)}function G(t){return function(){var e=this,r=arguments;return new Promise(function(n,a){var o=t.apply(e,r);function s(t){F(o,n,a,s,i,"next",t)}function i(t){F(o,n,a,s,i,"throw",t)}s(void 0)})}}var I={components:{LogoWide:w},data:function(){return{username:null,password:null,error:!1}},methods:{onSubmit:function(){var t=G(v.a.mark(function t(){var e,r,n=this;return v.a.wrap(function(t){while(1)switch(t.prev=t.next){case 0:return e=new URLSearchParams,e.append("username",this.username),e.append("password",this.password),t.next=5,this.$axios.post("/api/login",e).catch(function(t){t.response&&401===t.response.status?n.error="Login incorrect.":(console.error(t),n.error="Server error.")});case 5:if(r=t.sent,r){t.next=8;break}return t.abrupt("return");case 8:if(r.data.privScanTickets){t.next=11;break}return this.error="Not authorized to use ticket scanner.",t.abrupt("return");case 11:this.$emit("auth",r.data.token);case 12:case"end":return t.stop()}},t,this)}));function e(){return t.apply(this,arguments)}return e}()}},N=I,W=(r("537b"),Object(A["a"])(N,V,D,!1,null,null,null)),X=W.exports,U=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{attrs:{id:"scanner"}},[r("div",{attrs:{id:"scanner-top"}},[r("CameraView",{on:{decode:t.onDecode}}),r("div",{attrs:{id:"scanner-controls"}},[r("LogoTall"),r("form",{attrs:{id:"orderidform"},on:{submit:function(e){return e.preventDefault(),t.onSubmit(e)}}},[r("div",{attrs:{id:"forminner"}},[r("b-form-input",{attrs:{id:"orderid",type:"number",min:"1",step:"any",size:"5",placeholder:"Order #"},model:{value:t.orderid,callback:function(e){t.orderid=e},expression:"orderid"}}),r("b-button",{attrs:{type:"submit"}},[t._v("Submit")])],1),t.freeEntry?r("b-button",{attrs:{disabled:!t.freeEntry},on:{click:function(e){return e.preventDefault(),t.onFree(e)}}},[t._v("Free Entry")]):t._e()],1),r("div",{attrs:{id:"orderinfo"}},[r("div",{attrs:{id:"ordername"},domProps:{textContent:t._s(t.order&&t.order.name?t.order.name:" ")}}),r("div",{attrs:{id:"ordernum"},domProps:{textContent:t._s(t.order&&t.order.id?"Order number "+t.order.id:" ")}})])],1)],1),r("div",{attrs:{id:"scanner-bottom"}},[t.scanError?r("div",{attrs:{id:"scanError"},domProps:{textContent:t._s(t.scanError)}}):t.order&&t.order.classes?r("div",{attrs:{id:"quantities"}},t._l(t.order.classes,function(e){return r("ClassUsage",{key:e.name,attrs:{tclass:e},on:{change:function(r){return t.onCountChange(e,r)}}})}),1):t._e()])])},$=[],z=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{attrs:{id:"camera"}},[t.error?r("div",{attrs:{id:"error"},domProps:{textContent:t._s(t.error)}}):r("QrcodeStream",{attrs:{track:!1},on:{init:t.onStreamInit,decode:t.onDecode}})],1)},H=[],M=r("9a3e");function B(t,e,r,n,a,o,s){try{var i=t[o](s),c=i.value}catch(u){return void r(u)}i.done?e(c):Promise.resolve(c).then(n,a)}function J(t){return function(){var e=this,r=arguments;return new Promise(function(n,a){var o=t.apply(e,r);function s(t){B(o,n,a,s,i,"next",t)}function i(t){B(o,n,a,s,i,"throw",t)}s(void 0)})}}var Q={components:{QrcodeStream:M["QrcodeStream"]},data:function(){return{error:null}},computed:{streamStyle:function(){var t=3*(window.innerHeight-100)/4;t>window.innerWidth&&(t=window.innerWidth);var e=4*t/3;return{width:"".concat(t,"px"),height:"".concat(e,"px")}}},methods:{onDecode:function(t){this.$emit("decode",t)},onStreamInit:function(){var t=J(v.a.mark(function t(e){var r=this;return v.a.wrap(function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,e.catch(function(t){"NotAllowedError"===t.name?r.error="no permission to use camera":"NotFoundError"===t.name?r.error="no camera on this device":"NotSupportedError"===t.name?r.error="secure context required (HTTPS, localhost)":"NotReadableError"===t.name?r.error="camera already in use":"OverconstrainedError"===t.name?r.error="camera is not suitable":"StreamApiNotSupportedError"===t.name?r.error="browser can't use camera":r.error="can't start QR code scanner"});case 2:case"end":return t.stop()}},t)}));function e(e){return t.apply(this,arguments)}return e}()}},Y=Q,Z=(r("07bf"),Object(A["a"])(Y,z,H,!1,null,null,null)),K=Z.exports,q=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{staticClass:"class-usage",class:t.classColor},[r("div",{staticClass:"class-usage-buttons"},[r("div",{domProps:{textContent:t._s(t.tclass.name||"General Admission")}}),t._l(t.max,function(e){return r("b-btn",{key:e,staticClass:"class-usage-button",attrs:{variant:t.countVariant(e)},domProps:{textContent:t._s(e)},on:{click:function(r){return t.onClick(e)}}})})],2)])},tt=[],et={props:{tclass:Object},computed:{classColor:function(){return this.tclass.overflow?"class-usage-overflow":this.tclass.name?"class-usage-restricted":null},max:function(){return this.tclass.max<1e3?this.tclass.max:6*Math.ceil((this.tclass.used+1)/6)}},methods:{countVariant:function(t){return t<=this.tclass.min?"secondary":t<=this.tclass.used?"success":"outline-success"},onClick:function(t){t<this.tclass.min||(1===t&&0===this.tclass.min&&1===this.tclass.used?this.$emit("change",0):this.$emit("change",t))}}},rt=et,nt=(r("17b1"),Object(A["a"])(rt,q,tt,!1,null,null,null)),at=nt.exports,ot=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{attrs:{id:"logo-tall"}},[r("img",{attrs:{src:t.logo}}),t._v("Schola"),r("br"),t._v("Ticket"),r("br"),t._v("Scanner")])},st=[],it=r("8fd5"),ct=r.n(it),ut={computed:{logo:function(){return ct.a}}},lt=ut,dt=(r("474f"),Object(A["a"])(lt,ot,st,!1,null,null,null)),ft=dt.exports;function ht(t,e,r,n,a,o,s){try{var i=t[o](s),c=i.value}catch(u){return void r(u)}i.done?e(c):Promise.resolve(c).then(n,a)}function pt(t){return function(){var e=this,r=arguments;return new Promise(function(n,a){var o=t.apply(e,r);function s(t){ht(o,n,a,s,i,"next",t)}function i(t){ht(o,n,a,s,i,"throw",t)}s(void 0)})}}var vt={components:{CameraView:K,ClassUsage:at,LogoTall:ft},props:{auth:String,event:Object,freeEntry:Array},data:function(){return{order:null,orderid:null,scanError:null}},methods:{fetchTicket:function(){var t=pt(v.a.mark(function t(e){var r,n=this;return v.a.wrap(function(t){while(1)switch(t.prev=t.next){case 0:return this.scanError=this.order=null,t.next=3,this.$axios.post("/api/event/".concat(this.event.id,"/ticket/").concat(e),null,{headers:{auth:this.auth}}).catch(function(t){return t.response&&404===t.response.status?n.scanError="No such order":(console.log(t),n.scanError="Server error"),null});case 3:if(r=t.sent,r){t.next=6;break}return t.abrupt("return");case 6:r.data.id&&(this.order=r.data),r.data.error&&(this.scanError=r.data.error);case 8:case"end":return t.stop()}},t,this)}));function e(e){return t.apply(this,arguments)}return e}(),onCountChange:function(){var t=pt(v.a.mark(function t(e,r){var n,a,o=this;return v.a.wrap(function(t){while(1)switch(t.prev=t.next){case 0:return n=new URLSearchParams,n.append("scan",this.order.scan),n.append("class",e.name),n.append("used",r),t.next=6,this.$axios.post("/api/event/".concat(this.event.id,"/ticket/").concat(this.order.id||"free"),n,{headers:{auth:this.auth}}).catch(function(t){return console.log(t),o.scanError="Server error",null});case 6:if(a=t.sent,a){t.next=9;break}return t.abrupt("return");case 9:if(!a.data.error){t.next=12;break}return this.scanError=a.data.error,t.abrupt("return");case 12:this.$set(this.order,"scan",a.data.scan),this.$set(this.order,"id",a.data.id),e.used=r,e.overflow=!1;case 16:case"end":return t.stop()}},t,this)}));function e(e,r){return t.apply(this,arguments)}return e}(),onDecode:function(t){var e=t.match(/\/ticket\/(\d{4}-\d{4}-\d{4})$/);e?this.fetchTicket(e[1]):this.scanError="Not a Schola order"},onFree:function(){this.order={scan:"free",name:"Free Entry",classes:this.freeEntry.map(function(t){return{name:t,min:0,max:1e3,used:0}})}},onSubmit:function(){this.orderid>0&&this.fetchTicket(this.orderid),this.orderid=null}}},mt=vt,bt=(r("245e"),Object(A["a"])(mt,U,$,!1,null,null,null)),At=bt.exports,xt={components:{EventChooser:_,FatalError:R,LoginPage:X,TicketScanner:At},data:function(){return{auth:null,event:null,fatal:null}}},gt=xt,wt=(r("2906"),Object(A["a"])(gt,l,d,!1,null,null,null)),yt=wt.exports;n["default"].config.productionTip=!1,new n["default"]({render:function(t){return t(yt)}}).$mount("#main")},"58b5":function(t,e,r){"use strict";var n=r("4def"),a=r.n(n);a.a},5999:function(t,e,r){},"8fd5":function(t,e){t.exports="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADEAAAAaCAQAAAC0/nxiAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QA/4ePzL8AAAAHdElNRQfjBRkFMAkZpBhVAAADGklEQVRIx62VX2iVZRzHP2d/Wsvp0lIrmYmZhkuOYrsJo1ySWhRi0B+IAi+sLqS/QgVBVJB4sWStoJvdRBhSRCVkUGKWx4yT4WwTw3JWeGRbOcttTT3n08V5zrv31VGdHX/vzfv83uf5fn7P7/k9vxcqNEvPVFfaYbfnLFm/77jKuosDqHeNH3hC1REHHbNfXX8xAAvs9I8g2e0zZozboUrla7zTfUEs72cucb69CUShMsDlPu/xIPW3nV4rzrYrgfixEsAsOx0JQmd90ytFrPIV8xHgtC9NHNDsRzGp97w6qq0ZvmaXJ8yZ9SmnTkweb/ObWDK2OyecTJ11VlvrbJd5i01WSapcBAAr6GBB5PyCRxmkhUVczxXkGeAXusmSK34uCxEAq9nC/MiZoZ0m1rKYSzhEloP8xhkaOM5uzpW9AzHlgx4dK0j3+Lr7Lai9bnSe1VEqU5ZCKhPwkLnYGXznNk+p2uPtRdEyZc875FrXJQAH3BoAOe8uRV8Z4EkHYoCMG/0ppOvlYoImbBHgVAzwvct9N7wfcVEFOwiAyb7on/Gm4B22+HsYbfUyufApWc2/ywPQwAs8QX3k3sdz7GIDjWF8lOHwdg1LaQa6+JaB/3XpRGx0s6OxHexyiYhtkefVEPcK9zhiwYLDfuLc/0xdWDbTtz0TieX92IXhS3vk3STidWYTHXaTVUVE1XisyDeHLayjNriH6eAxesKol9KfYBaXAmluSIjcyOREtOPEj2l3xKLK+XTpWEFsjSrsoHPF+xxO7GK7U4KuTXFIJD/Nhz0QaxRfurrUHMK8Rj8MX0d9Vmx2fwKxOVL1LVutSpRbk4+4I1ak/baFZp3c6c3h6mmfjzvJxb7h5/5gQT1SLIri5Lv81DZXmrbZFu+13a7ogPP2+b6t1pyfzhDMPfaEmUNu835vdb0Z9Zhrx+anhDQPcBMNTGI606kGRjnNIIfZy1dkGYILu75F51I2sIoZAOQ5Sy1DfE0bO7G0JmXxAl7FPBYyk3pS/MVJjvEzOU6Wqma8SxSinMIylpNmGgX6OUyG3fTF1/wDOpCEvoS7FXkAAAAldEVYdGRhdGU6Y3JlYXRlADIwMTktMDUtMjVUMTI6NDg6MDktMDc6MDDEtiRYAAAAJXRFWHRkYXRlOm1vZGlmeQAyMDE5LTA1LTI1VDEyOjQ4OjA5LTA3OjAwteuc5AAAAABJRU5ErkJggg=="},"9c5c":function(t,e,r){},"9e51":function(t,e,r){},a4ff:function(t,e,r){"use strict";var n=r("9e51"),a=r.n(n);a.a},c3d6:function(t,e,r){"use strict";var n=r("d216"),a=r.n(n);a.a},d216:function(t,e,r){},f13f:function(t,e,r){}});
//# sourceMappingURL=app.baf714a8.js.map