(function(e){function t(t){for(var n,o,i=t[0],c=t[1],u=t[2],l=0,p=[];l<i.length;l++)o=i[l],Object.prototype.hasOwnProperty.call(a,o)&&a[o]&&p.push(a[o][0]),a[o]=0;for(n in c)Object.prototype.hasOwnProperty.call(c,n)&&(e[n]=c[n]);d&&d(t);while(p.length)p.shift()();return s.push.apply(s,u||[]),r()}function r(){for(var e,t=0;t<s.length;t++){for(var r=s[t],n=!0,i=1;i<r.length;i++){var c=r[i];0!==a[c]&&(n=!1)}n&&(s.splice(t--,1),e=o(o.s=r[0]))}return e}var n={},a={app:0},s=[];function o(t){if(n[t])return n[t].exports;var r=n[t]={i:t,l:!1,exports:{}};return e[t].call(r.exports,r,r.exports,o),r.l=!0,r.exports}o.m=e,o.c=n,o.d=function(e,t,r){o.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:r})},o.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},o.t=function(e,t){if(1&t&&(e=o(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var r=Object.create(null);if(o.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var n in e)o.d(r,n,function(t){return e[t]}.bind(null,n));return r},o.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return o.d(t,"a",t),t},o.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},o.p="";var i=window["webpackJsonp"]=window["webpackJsonp"]||[],c=i.push.bind(i);i.push=t,i=i.slice();for(var u=0;u<i.length;u++)t(i[u]);var d=c;s.push([0,"chunk-vendors"]),r()})({0:function(e,t,r){e.exports=r("56d7")},"0c8d":function(e,t,r){},"0ca4":function(e,t,r){},3122:function(e,t,r){},4713:function(e,t,r){"use strict";var n=r("0ca4"),a=r.n(n);a.a},"51af":function(e,t,r){},5283:function(e,t,r){},"56d7":function(e,t,r){"use strict";r.r(t);r("744f"),r("6c7b"),r("7514"),r("20d6"),r("1c4c"),r("6762"),r("cadf"),r("e804"),r("55dd"),r("d04f"),r("c8ce"),r("217b"),r("7f7f"),r("f400"),r("7f25"),r("536b"),r("d9ab"),r("f9ab"),r("32d7"),r("25c9"),r("9f3c"),r("042e"),r("c7c6"),r("f4ff"),r("049f"),r("7872"),r("a69f"),r("0b21"),r("6c1a"),r("c7c62"),r("84b4"),r("c5f6"),r("2e37"),r("fca0"),r("7cdf"),r("ee1d"),r("b1b1"),r("87f3"),r("9278"),r("5df2"),r("04ff"),r("f751"),r("4504"),r("fee7"),r("ffc1"),r("0d6d"),r("9986"),r("8e6e"),r("25db"),r("e4f7"),r("b9a1"),r("64d5"),r("9aea"),r("db97"),r("66c8"),r("57f0"),r("165b"),r("456d"),r("cf6a"),r("fd24"),r("8615"),r("551c"),r("097d"),r("df1b"),r("2397"),r("88ca"),r("ba16"),r("d185"),r("ebde"),r("2d34"),r("f6b3"),r("2251"),r("c698"),r("a19f"),r("9253"),r("9275"),r("3b2b"),r("3846"),r("4917"),r("a481"),r("28a5"),r("386d"),r("6b54"),r("4f7f"),r("8a81"),r("ac4d"),r("8449"),r("9c86"),r("fa83"),r("48c0"),r("a032"),r("aef6"),r("d263"),r("6c37"),r("9ec8"),r("5695"),r("2fdb"),r("d0b0"),r("5df3"),r("b54a"),r("f576"),r("ed50"),r("788d"),r("14b9"),r("f386"),r("f559"),r("1448"),r("673e"),r("242a"),r("c66f"),r("b05c"),r("34ef"),r("6aa2"),r("15ac"),r("af56"),r("b6e4"),r("9c29"),r("63d9"),r("4dda"),r("10ad"),r("c02b"),r("4795"),r("130f"),r("ac6a"),r("96cf"),r("0cdd");var n=r("2b0e"),a=r("bc3a"),s=r.n(a),o={},i=s.a.create(o);i.interceptors.request.use((function(e){return e}),(function(e){return Promise.reject(e)})),i.interceptors.response.use((function(e){return e}),(function(e){return Promise.reject(e)})),Plugin.install=function(e){e.axios=i,window.axios=i,Object.defineProperties(e.prototype,{axios:{get:function(){return i}},$axios:{get:function(){return i}}})},n["default"].use(Plugin);Plugin;var c=r("5f5b");r("ab8b"),r("2dd8");n["default"].use(c["a"]);var u=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{attrs:{id:"app"}},[e._m(0),r("div",{attrs:{id:"app-body"}},[e.$store.state.auth?r("Main"):r("Login")],1)])},d=[function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{attrs:{id:"app-header"}},[n("img",{attrs:{id:"app-header-logo",src:r("8fd5")}}),n("div",{attrs:{id:"app-header-title"}},[e._v("Schola Order Management")])])}],l=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("form",{attrs:{id:"login"},on:{submit:function(t){return t.preventDefault(),e.onSubmit(t)}}},[r("div",{attrs:{id:"login-header"}},[e._v("Please log in.")]),r("b-form-group",{attrs:{label:"Username","label-for":"login-username","label-cols":"4"}},[r("b-form-input",{attrs:{id:"login-username",autocapitalize:"none",autocomplete:"username",autofocus:"",trim:""},model:{value:e.username,callback:function(t){e.username=t},expression:"username"}})],1),r("b-form-group",{attrs:{label:"Password","label-for":"login-password","label-cols":"4"}},[r("b-form-input",{attrs:{id:"login-password",type:"password",autocomplete:"password"},model:{value:e.password,callback:function(t){e.password=t},expression:"password"}})],1),e.error?r("div",{attrs:{id:"login-error"},domProps:{textContent:e._s(e.error)}}):e._e(),r("b-button",{attrs:{id:"login-button",type:"submit",variant:"primary",disabled:e.disabled}},[e._v("Login")])],1)},p=[],f=r("a34a"),h=r.n(f);function m(e,t,r,n,a,s,o){try{var i=e[s](o),c=i.value}catch(u){return void r(u)}i.done?t(c):Promise.resolve(c).then(n,a)}function v(e){return function(){var t=this,r=arguments;return new Promise((function(n,a){var s=e.apply(t,r);function o(e){m(s,n,a,o,i,"next",e)}function i(e){m(s,n,a,o,i,"throw",e)}o(void 0)}))}}var b={data:function(){return{error:null,password:null,username:null}},computed:{disabled:function(){return!this.password||!this.username}},methods:{onSubmit:function(){var e=v(h.a.mark((function e(){var t,r;return h.a.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return t=new URLSearchParams,t.append("username",this.username),t.append("password",this.password),e.prev=3,e.next=6,this.$axios.post("/api/login",t);case 6:if(r=e.sent.data,r.privViewOrders){e.next=10;break}return this.error="Not authorized to use this app",e.abrupt("return");case 10:this.$store.commit("login",{auth:r.token,stripeKey:r.stripePublicKey,username:this.username}),e.next=21;break;case 13:if(e.prev=13,e.t0=e["catch"](3),!e.t0.response||401!==e.t0.response.status){e.next=18;break}return this.error="Login incorrect",e.abrupt("return");case 18:return console.error("Error logging in",e.t0),this.error="Server error — login failed",e.abrupt("return");case 21:case"end":return e.stop()}}),e,this,[[3,13]])})));function t(){return e.apply(this,arguments)}return t}()}},g=b,y=(r("b543"),r("2877")),C=Object(y["a"])(g,l,p,!1,null,null,null),k=C.exports,A=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("Report")},S=[],x=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("Split",{attrs:{id:"report"}},[r("SplitArea",{attrs:{size:10,minSize:200}},[e.report?r("ReportCriteria",{attrs:{stats:e.report},on:{update:e.onUpdate}}):r("div",{attrs:{id:"report-spinner"}},[r("b-spinner",{attrs:{label:"Loading..."}})],1)],1),r("SplitArea",{attrs:{size:90}},[e.report?e._e():r("div",{attrs:{id:"report-message"}}),r("div",{attrs:{id:"report-results"}},[e.haveParams?e.report.lines?e.report.lines.length?r("ReportTable",{attrs:{lines:e.report.lines}}):r("div",{attrs:{id:"report-message"}},[e._v("No orders match your search criteria.")]):r("div",{attrs:{id:"report-message"}},[e._v("Too many results; please narrow the search criteria.")]):r("div",{attrs:{id:"report-message"}},[e._v("For a list of orders, please provide search criteria.")]),r("div",{attrs:{id:"report-stats"},domProps:{textContent:e._s(e.stats)}})],1)])],1)},w=[],T=r("cba5"),E=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{attrs:{id:"report-criteria"}},[r("div",{attrs:{id:"report-criteria-heading"}},[e._v("Report Criteria")]),r("div",{attrs:{id:"report-criteria-scroll"}},[r("div",{staticClass:"report-criteria-section"},[e._v("Customer")]),r("input",{directives:[{name:"model",rawName:"v-model.lazy.trim",value:e.customer,expression:"customer",modifiers:{lazy:!0,trim:!0}}],attrs:{id:"report-criteria-customer",type:"text",placeholder:"Name or Email"},domProps:{value:e.customer},on:{change:function(t){e.customer=t.target.value.trim()},blur:function(t){return e.$forceUpdate()}}}),r("div",{staticClass:"report-criteria-section"},[e._v("Order Dates")]),r("div",{staticClass:"report-criteria-date-box"},[e._v("From"),r("input",{directives:[{name:"model",rawName:"v-model",value:e.createdAfter,expression:"createdAfter"}],staticClass:"report-criteria-date",attrs:{type:"date",placeholder:"From Date"},domProps:{value:e.createdAfter},on:{input:function(t){t.target.composing||(e.createdAfter=t.target.value)}}})]),r("div",{staticClass:"report-criteria-date-box"},[e._v("To"),r("input",{directives:[{name:"model",rawName:"v-model",value:e.createdBefore,expression:"createdBefore"}],staticClass:"report-criteria-date",attrs:{type:"date",placeholder:"To Date"},domProps:{value:e.createdBefore},on:{input:function(t){t.target.composing||(e.createdBefore=t.target.value)}}})]),r("div",{staticClass:"report-criteria-section"},[e._v("Products")]),r("TreeSelect",{attrs:{tree:e.productsTree,value:e.selectedProducts},on:{change:e.onChangeProducts}}),r("div",{staticClass:"report-criteria-section"},[e._v("Tickets Used At")]),r("TreeSelect",{attrs:{tree:e.usedAtEventsTree,value:e.selectedUsedAtEvents},on:{change:e.onChangeUsedAtEvents}}),r("div",{staticClass:"report-criteria-section"},[e._v("Ticket Classes")]),r("TreeSelect",{attrs:{tree:e.ticketClassList,value:e.selectedTicketClasses},on:{change:e.onChangeTicketClasses}}),r("div",{staticClass:"report-criteria-section"},[e._v("Order Sources")]),r("TreeSelect",{attrs:{tree:e.orderSourcesList,value:e.selectedOrderSources},on:{change:e.onChangeOrderSources}}),r("div",{staticClass:"report-criteria-section"},[e._v("Payment Types")]),r("TreeSelect",{attrs:{tree:e.paymentTypesTree,value:e.selectedPaymentTypes},on:{change:e.onChangePaymentTypes}}),r("div",{staticClass:"report-criteria-section"},[e._v("Coupon Codes")]),r("TreeSelect",{attrs:{tree:e.orderCouponsList,value:e.selectedOrderCoupons},on:{change:e.onChangeOrderCoupons}})],1),r("div",{attrs:{id:"report-criteria-buttons"}},[r("b-button",{attrs:{variant:"outline-primary"},on:{click:e.onReset}},[e._v("Reset")]),r("b-button",{attrs:{variant:"primary"},on:{click:e.onSearch}},[e._v("Search")])],1)])},O=[],P=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"tree-select"},[e._l(e.tree,(function(t){return[t.children?r("TreeSelectAncestor",{key:t.aid,attrs:{checked:e.checked,count:e.count,expanded:e.expanded,node:t,path:[t.aid],value:e.value},on:{bulkChange:function(t){return e.onBulkChange(t)},change:function(t){return e.onChange(t)},expand:function(t){return e.onExpand(t)}}}):r("TreeSelectLeaf",{key:t.id,attrs:{flat:e.flat,node:t,path:[],value:e.value.includes(t.id)},on:{change:function(t){return e.onChange(t)}}})]}))],2)},j=[],B=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"tree-select-ancestor"},[r("div",{staticClass:"tree-select-item"},[r("Icon",{staticClass:"tree-select-expand",attrs:{"fixed-width":"",icon:e.isExpanded?e.caretDown:e.caretRight},on:{click:e.onExpand}}),r("input",{ref:"cb",staticClass:"tree-select-cb",attrs:{type:"checkbox"},domProps:{checked:e.isChecked},on:{change:e.onChange}}),r("div",{staticClass:"tree-select-label",domProps:{textContent:e._s(e.node.label)}}),r("div",{staticClass:"tree-select-count",domProps:{textContent:e._s(e.node.count)}})],1),e.isExpanded?r("div",{staticClass:"tree-select-ancestor-indent"},[e._l(e.node.children,(function(t){return[t.children?r("TreeSelectAncestor",{key:t.label,attrs:{checked:e.checked,count:e.count,expanded:e.expanded,node:t,path:e.path.concat([t.aid]),value:e.value},on:{bulkChange:function(t){return e.$emit("bulkChange",t)},change:function(t){return e.$emit("change",t)},expand:function(t){return e.$emit("expand",t)}}}):r("TreeSelectLeaf",{key:t.id,attrs:{flat:!1,node:t,path:e.path,value:e.value.includes(t.id)},on:{change:function(t){return e.$emit("change",t)}}})]}))],2):e._e()])},_=[],J=r("ad3d"),L=r("c074"),U=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"tree-select-item"},[e.flat?e._e():r("div",{staticClass:"tree-select-leaf-expand"}),r("input",{staticClass:"tree-select-cb",attrs:{type:"checkbox"},domProps:{checked:e.value},on:{change:e.onChange}}),r("div",{staticClass:"tree-select-label",domProps:{textContent:e._s(e.node.label||e.node.id)}}),r("div",{staticClass:"tree-select-count",domProps:{textContent:e._s(e.node.count)}})])},R=[],M={props:{flat:Boolean,node:Object,path:Array,value:Boolean},methods:{onChange:function(){this.$emit("change",{id:this.node.id,path:this.path,checked:!this.value})}}},Q=M,D=(r("f0b2"),Object(y["a"])(Q,U,R,!1,null,null,null)),I=D.exports,Z={name:"TreeSelectAncestor",components:{Icon:J["a"],TreeSelectLeaf:I},props:{checked:Object,count:Object,expanded:Object,node:Object,path:Array,value:Array},data:function(){return{caretDown:L["a"],caretRight:L["b"]}},computed:{isChecked:function(){return this.checked[this.node.aid]===this.count[this.node.aid]},isExpanded:function(){return this.expanded[this.node.aid]},isIndeterminate:function(){return!this.isChecked&&0!==this.checked[this.node.aid]}},watch:{isIndeterminate:{immediate:!0,handler:function(){var e=this;this.$refs.cb?this.$refs.cb.indeterminate=this.isIndeterminate:this.$nextTick((function(){e.$refs.cb.indeterminate=e.isIndeterminate}))}}},methods:{onChange:function(){this.$emit("bulkChange",{node:this.node,path:this.path,checked:!this.isChecked})},onExpand:function(){this.$emit("expand",this.node.aid)}}},z=Z,K=(r("6975"),Object(y["a"])(z,B,_,!1,null,null,null)),G=K.exports;function N(e){return W(e)||q(e)||$()}function $(){throw new TypeError("Invalid attempt to spread non-iterable instance")}function q(e){if(Symbol.iterator in Object(e)||"[object Arguments]"===Object.prototype.toString.call(e))return Array.from(e)}function W(e){if(Array.isArray(e)){for(var t=0,r=new Array(e.length);t<e.length;t++)r[t]=e[t];return r}}var V={components:{TreeSelectAncestor:G,TreeSelectLeaf:I},props:{tree:Array,value:Array},data:function(){return{count:{},checked:{},expanded:{},flat:!1}},watch:{tree:{immediate:!0,handler:"setup"}},methods:{onBulkChange:function(e){var t=this,r=e.node,n=e.path,a=e.checked,s=N(this.value);r.children.forEach((function(e){s=t.onBulkChangeInner(e,n,a,s)})),this.$emit("change",s)},onBulkChangeInner:function(e,t,r,n){var a=this;return e.children?(t=[].concat(N(t),[e.aid]),e.children.forEach((function(e){n=a.onBulkChangeInner(e,t,r,n)}))):r&&!n.includes(e.id)?(n.push(e.id),t.forEach((function(e){a.checked[e]++}))):!r&&n.includes(e.id)&&(n=n.filter((function(t){return t!==e.id})),t.forEach((function(e){a.checked[e]--}))),n},onChange:function(e){var t=this,r=e.id,n=e.checked,a=e.path;n?(this.$emit("change",[r].concat(N(this.value))),a.forEach((function(e){t.checked[e]++}))):(this.$emit("change",this.value.filter((function(e){return e!==r}))),a.forEach((function(e){t.checked[e]--})))},onExpand:function(e){this.expanded[e]=!this.expanded[e]},setup:function(){var e=this;this.count={},this.checked={},this.flat=!0,this.tree.forEach((function(t){e.setupInner([],t)}))},setupInner:function(e,t){var r=this;t.children?(this.flat=!1,e=[].concat(N(e),[t.aid]),this.$set(this.count,t.aid,0),this.$set(this.checked,t.aid,0),this.$set(this.expanded,t.aid,!!this.expanded[t.aid]),t.children.forEach((function(t){r.setupInner(e,t)}))):(e.forEach((function(e){r.count[e]++})),this.value.includes(t.id)&&e.forEach((function(e){r.checked[e]++})))}}},H=V,F=Object(y["a"])(H,P,j,!1,null,null,null),Y=F.exports,X={auctionitem:"Auction Items",donation:"Donations",recording:"Recordings",sheetmusic:"Sheet Music",ticket:"Tickets",other:"[other]"},ee={gala:"Gala Software",inperson:"Event Box Office",members:"Members Web Site",office:"Schola Office",public:"Public Web Site"},te={components:{TreeSelect:Y},props:{stats:Object},data:function(){return{createdAfter:"",createdBefore:"",customer:"",selectedOrderCoupons:[],selectedOrderSources:[],selectedPaymentTypes:[],selectedProducts:[],selectedTicketClasses:[],selectedUsedAtEvents:[],updateTimer:null}},computed:{orderCouponsList:function(){return this.stats.orderCoupons.map((function(e){return{id:e.n,label:e.n||"(none)",count:e.c}})).sort((function(e,t){return e.label.localeCompare(t.label)}))},orderSourcesList:function(){return this.stats.orderSources.map((function(e){return{id:e.os,label:ee[e.os]||e.os,count:e.c}})).sort((function(e,t){return e.label.localeCompare(t.label)}))},paymentTypesTree:function(){var e={};return this.stats.paymentTypes.forEach((function(t){var r=t.n.split(","),n=r.pop(),a=e,s="";r.forEach((function(e){s+=","+e,a[e]||(a[e]={label:e,aid:s,count:0,children:{}}),a[e].count+=t.c,a=a[e].children})),a[t.n]={id:t.n,label:n,count:t.c}})),e=this.sortTree(e),e},productsTree:function(){var e={};return this.stats.products.forEach((function(t){e[t.ptype]||(e[t.ptype]={label:X[t.ptype]||t.ptype,aid:t.ptype,count:0,children:{}}),t.series?(e[t.ptype].children[t.series]||(e[t.ptype].children[t.series]={label:t.series,aid:"".concat(t.ptype," ").concat(t.series),count:0,children:{}}),e[t.ptype].children[t.series].children[t.id]={id:t.id,label:t.name,count:t.count},e[t.ptype].children[t.series].count+=t.count):e[t.ptype].children[t.id]={id:t.id,label:t.name,count:t.count},e[t.ptype].count+=t.count})),e=this.sortTree(e),e},ticketClassList:function(){return this.stats.ticketClasses.map((function(e){return{id:e.n,label:e.n||"General Admission",count:e.c}}))},usedAtEventsTree:function(){var e={};return this.stats.usedAtEvents.forEach((function(t){t.series?(e[t.series]||(e[t.series]={label:t.series,aid:t.series,count:0,children:{}}),e[t.series].children[t.id]={id:t.id,label:"".concat(t.start.substr(0,10)," ").concat(t.name),count:t.count},e[t.series].count+=t.count):e[t.id]={id:t.id,label:"".concat(t.start?t.start.substr(0,10)+" ":"").concat(t.name),count:t.count}})),e=this.sortTree(e),e}},methods:{onChangeOrderCoupons:function(e){this.selectedOrderCoupons=e},onChangeOrderSources:function(e){this.selectedOrderSources=e},onChangePaymentTypes:function(e){this.selectedPaymentTypes=e},onChangeProducts:function(e){this.selectedProducts=e},onChangeTicketClasses:function(e){this.selectedTicketClasses=e},onChangeUsedAtEvents:function(e){this.selectedUsedAtEvents=e},onReset:function(){this.customer="",this.createdAfter="",this.createdBefore="",this.selectedOrderCoupons=[],this.selectedOrderSources=[],this.selectedPaymentTypes=[],this.selectedProducts=[],this.selectedTicketClasses=[],this.selectedUsedAtEvents=[],this.onSearch()},onSearch:function(){var e=new URLSearchParams;this.customer&&e.append("customer",this.customer),this.createdAfter&&e.append("createdAfter",this.createdAfter+"T00:00:00"),this.createdBefore&&e.append("createdBefore",this.createdBefore+"T23:59:59"),this.selectedOrderCoupons.forEach((function(t){e.append("orderCoupon",t)})),this.selectedOrderSources.forEach((function(t){e.append("orderSource",t)})),this.selectedPaymentTypes.forEach((function(t){e.append("paymentType",t)})),this.selectedProducts.forEach((function(t){e.append("product",t)})),this.selectedTicketClasses.forEach((function(t){e.append("ticketClass",t)})),this.selectedUsedAtEvents.forEach((function(t){e.append("usedAtEvent",t)})),this.$emit("update",e)},sortTree:function(e){for(var t in e)e[t].children&&(e[t].children=this.sortTree(e[t].children));return Object.keys(e).sort((function(t,r){return(e[t].id||e[t].label).localeCompare(e[r].id||e[r].label)})).map((function(t){return e[t]}))}}},re=te,ne=(r("d0f7"),Object(y["a"])(re,E,O,!1,null,null,null)),ae=ne.exports,se=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{attrs:{id:"report-table"}},[r("b-table",{attrs:{borderless:"",fields:e.fields,items:e.lines,small:"",striped:""}})],1)},oe=[],ie={gala:"Gala",inperson:"Box Office",members:"Members Site",office:"Schola Office",public:"Web Site"},ce={props:{lines:Array},data:function(){return{fields:[{key:"orderID",label:"Order",sortable:!0,sortdirection:"desc"},{key:"orderTime",label:"Date/Time",sortable:!0,sortdirection:"desc",formatter:function(e){return e.substr(0,16).replace("T"," ")}},{key:"name",label:"Customer",sortable:!0},{key:"email",label:"Email",sortable:!0},{key:"quantity",label:"Qty",sortable:!0,tdClass:"report-table-right"},{key:"product",label:"Product",sortable:!0},{key:"usedAtEvent",label:"Used At",sortable:!0,formatter:function(e){return e||"(unused)"}},{key:"orderSource",label:"Source",sortable:!0,formatter:function(e){return ie[e]||e}},{key:"paymentType",label:"Payment",sortable:!0,formatter:function(e){return e.replace(",",", ")}},{key:"amount",label:"Amount",sortable:!0,tdClass:"report-table-right",formatter:function(e){return"$".concat(e.toFixed(2))}}]}}},ue=ce,de=(r("fda5"),Object(y["a"])(ue,se,oe,!1,null,null,null)),le=de.exports;function pe(e,t,r,n,a,s,o){try{var i=e[s](o),c=i.value}catch(u){return void r(u)}i.done?t(c):Promise.resolve(c).then(n,a)}function fe(e){return function(){var t=this,r=arguments;return new Promise((function(n,a){var s=e.apply(t,r);function o(e){pe(s,n,a,o,i,"next",e)}function i(e){pe(s,n,a,o,i,"throw",e)}o(void 0)}))}}var he={components:{ReportCriteria:ae,ReportTable:le,Split:T["a"],SplitArea:T["a"].SplitArea},data:function(){return{haveParams:!1,report:null}},mounted:function(){this.runReport(null)},computed:{stats:function(){return"Matched ".concat(this.report.orderCount," ").concat(1===this.report.orderCount?"order":"orders",", ").concat(this.report.itemCount," ").concat(1===this.report.itemCount?"item":"items",", total amount $").concat(this.report.totalAmount.toFixed(2),".")}},methods:{onUpdate:function(e){this.haveParams=""!==e.toString(),this.runReport(e)},runReport:function(){var e=fe(h.a.mark((function e(t){return h.a.wrap((function(e){while(1)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,this.$axios.get("/api/report",{headers:{Auth:this.$store.state.auth},params:t});case 3:this.report=e.sent.data,e.next=9;break;case 6:e.prev=6,e.t0=e["catch"](0),window.alert(e.t0.toString());case 9:case"end":return e.stop()}}),e,this,[[0,6]])})));function t(t){return e.apply(this,arguments)}return t}()}},me=he,ve=(r("4713"),Object(y["a"])(me,x,w,!1,null,null,null)),be=ve.exports,ge={components:{Report:be}},ye=ge,Ce=Object(y["a"])(ye,A,S,!1,null,null,null),ke=Ce.exports,Ae={components:{Login:k,Main:ke},mounted:function(){var e=.01*window.innerHeight;document.documentElement.style.setProperty("--vh","".concat(e,"px"))}},Se=Ae,xe=(r("7faf"),Object(y["a"])(Se,u,d,!1,null,null,null)),we=xe.exports,Te=r("2f62");n["default"].use(Te["a"]);var Ee=new Te["a"].Store({state:{auth:null,stripeKey:null,username:null},mutations:{login:function(e,t){var r=t.auth,n=t.stripeKey,a=t.username;e.auth=r,e.stripeKey=n,e.username=a},logout:function(e){e.auth=e.stripeKey=e.username=null}}}),Oe=Ee;n["default"].config.productionTip=!1,new n["default"]({store:Oe,render:function(e){return e(we)}}).$mount("#app")},6975:function(e,t,r){"use strict";var n=r("3122"),a=r.n(n);a.a},"7faf":function(e,t,r){"use strict";var n=r("8fba"),a=r.n(n);a.a},"8fba":function(e,t,r){},"8fd5":function(e,t){e.exports="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAFoAAAAwCAYAAACRx20+AAAF5klEQVR42u2bbWiVZRjH//e2NNPmfEETfKHULEsTNctFmpQyK1d+MC2rD1mkQZgKYlEQQmVlkvZiWRQuFCdWBtqLoBLqh3QqOlvktGLzLaazMufL3H592DUcp7Pz3M+zOec554LD2LP7vp/r+e167vt6O1Ja/ifA1UAX4ElgJXAIqCVYyoBlwDCgHZCRphkfcC8gD/iapss+YDLQLk32IuCewBxgN80vr6Vh10GeDuwFajygbbVtJIzUAhNSGfAQYBtQ5QHrZ2AskAtUR7DqT1IRsAMeASo9rXE5kG1zC6PuH6kGuTuwwJPNceCVmPk/pUEHQ74d+NGTy9F4+yqwJSLnY6l04JV5bBMA24HRjazzQUTQX6RC4LEIOO0JZAtwY4L1BkQE/VQyQ24PfBcCxm6gb8CaWcCKkJC/BbokK+ROIffTcmBkiLU3AOcD1qwyyL2TFfJw4GDMvhvkXeTXu36e98gCngV+AH4BDtv2VGG/bwaeS+bt4nF76DBR26QwkOMAvwW4B5gCjAcGAlnJCjgDmAWcDAl5SkvpmCzkZ0t6O8T485KmSSqM80+7SVKupG6ShkpqJ6k+4EDSTkl/SjrgnNuYSoHIrAhJniX22jvzTnIt1E50wFUD54Cz9rM+57ENeDRpLdqS6nMlvWGW5rPHIuljSS9Iul7SGEnzJMVz62okHZZUbla8Q1KlpBJJmZJqzer7S8oFdjnnfk02K24LzI8QOLwJdAZeAo4kGLcHmAkMT+X0Znvgc+BCSMgzgfwA//o0MDtRdJhKoNdEsOQZwKvAvwnGlAIjorh5yQg5Si1vPjAv4KA7BPRLA67bk7+KAHk98GJAUqkSGJKGXJesXxPSfQNYZxXt/QHjJ0eNDJMJct8QyfqGkAs95y5LWzJ0NTcrLOQCoIPty4nkr6C0aKq4cAci7MmFQDbQwyON+V5K91uYJR+MAHkD0NnWWBgw9hxwbypDvgMoiQC5uGEU59GrUZSy2wYwxiodYffk7cDgBuuM9Ygal7d0A2JWa7FkSSslXec7xZJIBZJeds6VN/hbnkeC6aRzrtZTtxskDbLU6VWSKiRtv+JSpMCtIarUDWUFkBNnvaUeJazFHno9bG/LWeCMvSUX7JA9Y0WGhVcK5MERkkM1wPtA20bWLPZYY3GAXgUh9ClqzYCd9Q6HteTzwIcBa+/1WOfdBPM/ivB2vdVaQT/j2cUZC3lOUKhsFY9I3Z3AnVbJDitHgK6J9Mq4DJDnSlqkulqcr5RLynfOveOcC2oY/F0Xa3yNSSfgmjjX8yV1jfBY7SUNbzWggaWSFkjqEGJakaQ859z3nuO3qq7MlEj6S+oV53p2Ex6vzWUHbYXQ1ZKmy6+2V+/CrZV0t3OuJMTt9ntY9CBJ/Rq55yWRjBaAfIOkTZImhZhWKWmJc26ic+5smPs55zZJKgsaJmkUEGuF1VEfU1JVEIiMSwQ4E3gsQnKoCMhrhsM2SCqAgTHzxgOnIhyGB32U+vQSQB4HrArpvtVY/3HPZrh/BrDDs003J2be8xFAv+6jVAnwZf33NJr4gKOs8e94SEUrganN2bMG3OwZDO1q6IEY7KHAak/d/wE6+ih0m034I8orC+RYt1Ax0eSzxqK8ZoD9gKe/XgGMi/NWZNqXPLsDExuZOyKMQhMsfq//Ll0e0M16gTsA19on264PAKZZDjiKnLL79GmBw/iJEEHIZjtXettz9rB2sbWNZA8fiqJQHrAxZrGdVn3+xj7rQ3ZsNpST1lVfANzVwv77fRG+6HMiwd9+A+4P5Q3FKNRR0oOSRkoaK6mpHTt/S9qjur61rZJ2OefKdBnEQuQZkp6WFLULv1rScnM9iyODbqBUG1Omuzn3IyQNkZRjOePY8LnUIqNSC5f3SSqWdEzSCUlHPULnFkvLWs56qj2Tj5RIWidplaT9zrnTYe/rPBTLsAJBlo3PiAl0UF3npSz0rZFU45yracXZw0wzjBxJwySNNmvtY4ZRJem46rpIS+3aBedc1IBG/wEuWf8cvwSJcwAAAABJRU5ErkJggg=="},"996d":function(e,t,r){},b543:function(e,t,r){"use strict";var n=r("996d"),a=r.n(n);a.a},d0f7:function(e,t,r){"use strict";var n=r("5283"),a=r.n(n);a.a},f0b2:function(e,t,r){"use strict";var n=r("51af"),a=r.n(n);a.a},fda5:function(e,t,r){"use strict";var n=r("0c8d"),a=r.n(n);a.a}});
//# sourceMappingURL=app.fedfd682.js.map