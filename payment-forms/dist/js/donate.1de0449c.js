(function(t){function e(e){for(var a,s,o=e[0],u=e[1],c=e[2],d=0,m=[];d<o.length;d++)s=o[d],Object.prototype.hasOwnProperty.call(i,s)&&i[s]&&m.push(i[s][0]),i[s]=0;for(a in u)Object.prototype.hasOwnProperty.call(u,a)&&(t[a]=u[a]);l&&l(e);while(m.length)m.shift()();return r.push.apply(r,c||[]),n()}function n(){for(var t,e=0;e<r.length;e++){for(var n=r[e],a=!0,o=1;o<n.length;o++){var u=n[o];0!==i[u]&&(a=!1)}a&&(r.splice(e--,1),t=s(s.s=n[0]))}return t}var a={},i={donate:0},r=[];function s(e){if(a[e])return a[e].exports;var n=a[e]={i:e,l:!1,exports:{}};return t[e].call(n.exports,n,n.exports,s),n.l=!0,n.exports}s.m=t,s.c=a,s.d=function(t,e,n){s.o(t,e)||Object.defineProperty(t,e,{enumerable:!0,get:n})},s.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},s.t=function(t,e){if(1&e&&(t=s(t)),8&e)return t;if(4&e&&"object"===typeof t&&t&&t.__esModule)return t;var n=Object.create(null);if(s.r(n),Object.defineProperty(n,"default",{enumerable:!0,value:t}),2&e&&"string"!=typeof t)for(var a in t)s.d(n,a,function(e){return t[e]}.bind(null,a));return n},s.n=function(t){var e=t&&t.__esModule?function(){return t["default"]}:function(){return t};return s.d(e,"a",e),e},s.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},s.p="/";var o=window["webpackJsonp"]=window["webpackJsonp"]||[],u=o.push.bind(o);o.push=e,o=o.slice();for(var c=0;c<o.length;c++)e(o[c]);var l=u;r.push([1,"chunk-vendors"]),n()})({1:function(t,e,n){t.exports=n("ec2b")},"135e":function(t,e,n){"use strict";var a=n("2b0e"),i=n("5f5b");n("2dd8");a["default"].use(i["a"])},"1b96":function(t,e,n){"use strict";var a=n("3537"),i=n.n(a);i.a},3537:function(t,e,n){},"5d91":function(t,e,n){},"88a8":function(t,e,n){},"98b7":function(t,e,n){"use strict";var a=n("5d91"),i=n.n(a);i.a},"9a15":function(t,e,n){"use strict";var a=n("88a8"),i=n.n(a);i.a},a07c:function(t,e,n){},be3b:function(t,e,n){"use strict";var a=n("2b0e"),i=n("bc3a"),r=n.n(i),s={},o=r.a.create(s);o.interceptors.request.use((function(t){return t}),(function(t){return Promise.reject(t)})),o.interceptors.response.use((function(t){return t}),(function(t){return Promise.reject(t)})),Plugin.install=function(t){t.axios=o,window.axios=o,Object.defineProperties(t.prototype,{axios:{get:function(){return o}},$axios:{get:function(){return o}}})},a["default"].use(Plugin);Plugin},c798:function(t,e,n){"use strict";var a=n("a07c"),i=n.n(a);i.a},ec2b:function(t,e,n){"use strict";n.r(e);n("744f"),n("6095"),n("6c7b"),n("d25f"),n("7514"),n("20d6"),n("f3e2"),n("1c4c"),n("6762"),n("57e7"),n("2caf"),n("cadf"),n("9865"),n("6d67"),n("e804"),n("0cd8"),n("48f8"),n("759f"),n("55dd"),n("d04f"),n("78ce"),n("8ea5"),n("0298"),n("c8ce"),n("87b3"),n("d92a"),n("217b"),n("7f7f"),n("f400"),n("7f25"),n("536b"),n("d9ab"),n("f9ab"),n("32d7"),n("25c9"),n("9f3c"),n("042e"),n("c7c6"),n("f4ff"),n("049f"),n("7872"),n("a69f"),n("0b21"),n("6c1a"),n("c7c62"),n("84b4"),n("c5f6"),n("2e37"),n("fca0"),n("7cdf"),n("ee1d"),n("b1b1"),n("87f3"),n("9278"),n("5df2"),n("04ff"),n("f751"),n("8478"),n("4504"),n("fee7"),n("1c01"),n("58b2"),n("ffc1"),n("0d6d"),n("9986"),n("8e6e"),n("25db"),n("e4f7"),n("b9a1"),n("64d5"),n("9aea"),n("db97"),n("66c8"),n("57f0"),n("165b"),n("456d"),n("cf6a"),n("fd24"),n("8615"),n("551c"),n("097d"),n("df1b"),n("2397"),n("88ca"),n("ba16"),n("d185"),n("ebde"),n("2d34"),n("f6b3"),n("2251"),n("c698"),n("a19f"),n("9253"),n("9275"),n("3b2b"),n("3846"),n("4917"),n("a481"),n("28a5"),n("386d"),n("6b54"),n("4f7f"),n("8a81"),n("ac4d"),n("8449"),n("9c86"),n("fa83"),n("48c0"),n("a032"),n("aef6"),n("d263"),n("6c37"),n("9ec8"),n("5695"),n("2fdb"),n("d0b0"),n("5df3"),n("b54a"),n("f576"),n("ed50"),n("788d"),n("14b9"),n("f386"),n("f559"),n("1448"),n("673e"),n("242a"),n("4f37"),n("c66f"),n("262f"),n("b05c"),n("34ef"),n("6aa2"),n("15ac"),n("af56"),n("b6e4"),n("9c29"),n("63d9"),n("4dda"),n("10ad"),n("c02b"),n("4795"),n("130f"),n("ac6a"),n("96cf"),n("0cdd");var a=n("2b0e"),i=(n("be3b"),n("135e"),function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{attrs:{id:"donate-div"}},[n("Dialog",{ref:"dialog",attrs:{ordersURL:t.ordersURL,stripeKey:t.stripeKey}}),n("b-btn",{attrs:{id:"donate-btn",variant:"primary"},on:{click:function(e){return t.$refs.dialog.show()}}},[t._v("Donate Online")]),n("span",[t._v("with a credit or debit card")])],1)}),r=[],s=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("b-modal",{ref:"modal",attrs:{title:"Donation","no-close-on-backdrop":"","hide-footer":""},on:{shown:t.onShown,hide:t.onHide}},[t.orderID?n("Confirmation",{key:t.seq,attrs:{orderID:t.orderID},on:{close:t.onClose}}):n("OrderForm",{key:t.seq,ref:"form",attrs:{ordersURL:t.ordersURL,stripeKey:t.stripeKey},on:{success:t.onOrderSuccess,cancel:t.onClose,submitting:t.onSubmitting}})],1)},o=[],u=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",[n("div",{attrs:{id:"donate-confirm"}},[t._v("We have received your donation and emailed you a receipt.  We will send\na confirmation by postal mail for your tax records.  Thank you for\nsupporting Schola Cantorum!")]),n("div",{attrs:{id:"connect-head"}},[t._v("Stay informed of Schola news!")]),n("div",[n("b-spinner",{directives:[{name:"show",rawName:"v-show",value:t.emailSpinner,expression:"emailSpinner"}],staticClass:"connect-done mt-1",attrs:{small:""}}),n("div",{directives:[{name:"show",rawName:"v-show",value:t.emailDone,expression:"emailDone"}],staticClass:"connect-done"},[t._v("✓")]),n("div",{staticClass:"connect-link",on:{click:t.onEmail}},[t._v("Subscribe to our email list")])],1),n("div",[n("b-spinner",{directives:[{name:"show",rawName:"v-show",value:t.pmailSpinner,expression:"pmailSpinner"}],staticClass:"connect-done mt-1",attrs:{small:""}}),n("div",{directives:[{name:"show",rawName:"v-show",value:t.pmailDone,expression:"pmailDone"}],staticClass:"connect-done"},[t._v("✓")]),n("div",{staticClass:"connect-link",on:{click:t.onPMail}},[t._v("Subscribe to our postal mailing list")])],1),n("div",[n("div",{directives:[{name:"show",rawName:"v-show",value:t.facebookDone,expression:"facebookDone"}],staticClass:"connect-done"},[t._v("✓")]),n("div",{staticClass:"connect-link",on:{click:t.onFacebook}},[t._v("Follow us on Facebook")])]),n("div",[n("div",{directives:[{name:"show",rawName:"v-show",value:t.twitterDone,expression:"twitterDone"}],staticClass:"connect-done"},[t._v("✓")]),n("div",{staticClass:"connect-link",on:{click:t.onTwitter}},[t._v("Follow us on Twitter")])]),n("div",{staticStyle:{"text-align":"right"}},[n("b-btn",{attrs:{id:"connect-close",variant:"primary"},on:{click:function(e){return t.$emit("close")}}},[t._v("Close")])],1)])},c=[],l=n("a34a"),d=n.n(l);function m(t,e,n,a,i,r,s){try{var o=t[r](s),u=o.value}catch(c){return void n(c)}o.done?e(u):Promise.resolve(u).then(a,i)}function p(t){return function(){var e=this,n=arguments;return new Promise((function(a,i){var r=t.apply(e,n);function s(t){m(r,a,i,s,o,"next",t)}function o(t){m(r,a,i,s,o,"throw",t)}s(void 0)}))}}var f,h={props:{orderID:Number},data:function(){return{emailDone:!1,emailSpinner:!1,facebookDone:!1,pmailDone:!1,pmailSpinner:!1,twitterDone:!1}},methods:{onEmail:function(){var t=p(d.a.mark((function t(){var e;return d.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:if(!this.emailDone){t.next=2;break}return t.abrupt("return");case 2:return this.emailSpinner=!0,t.next=5,this.$axios.post("/backend/email-signup?order=".concat(this.orderID)).catch((function(t){return console.error(t),null}));case 5:if(e=t.sent,this.emailSpinner=!1,e){t.next=9;break}return t.abrupt("return");case 9:e.status<400?this.emailDone=!0:console.error(e.statusText);case 10:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),onFacebook:function(){window.open("https://www.facebook.com/scholacantorum.org","_blank"),this.facebookDone=!0},onPMail:function(){var t=p(d.a.mark((function t(){var e;return d.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:if(!this.pmailDone){t.next=2;break}return t.abrupt("return");case 2:return this.pmailSpinner=!0,t.next=5,this.$axios.post("/backend/mail-signup?order=".concat(this.orderID)).catch((function(t){return console.error(t),null}));case 5:if(e=t.sent,this.pmailSpinner=!1,e){t.next=9;break}return t.abrupt("return");case 9:e.status<400?this.pmailDone=!0:console.error(e.statusText);case 10:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),onTwitter:function(){window.open("https://twitter.com/scholacantorum1","_blank"),this.twitterDone=!0}}},b=h,v=(n("1b96"),n("2877")),y=Object(v["a"])(b,u,c,!1,null,null,null),g=y.exports,w=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("b-form",{attrs:{novalidate:""},on:{submit:function(e){return e.preventDefault(),t.onSubmit(e)}}},[n("b-form-group",{attrs:{state:t.amountState,"invalid-feedback":"Please enter an amount."}},[n("table",{attrs:{id:"donate-amount-row"}},[n("tr",[n("td",{attrs:{id:"donate-amount-label"}},[n("label",{staticClass:"mb-0",attrs:{for:"donate-amount"}},[t._v("Donation amount?")])]),n("td",{attrs:{id:"donate-amount-cell"}},[n("b-input-group",{attrs:{prepend:"$"}},[n("b-form-input",{ref:"amount",attrs:{id:"donate-amount",value:t.amount||"",disabled:t.submitting,type:"number",placeholder:"0",min:"0"},on:{input:function(e){t.amount=Math.max(parseInt(e)||0,0)}}})],1)],1)])])]),n("OrderPayment",{ref:"pmt",attrs:{send:t.onSend,stripeKey:t.stripeKey,total:100*t.amount||null},on:{cancel:t.onCancel,submitted:t.onSubmitted,submitting:t.onSubmitting}})],1)},x=[],S=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{directives:[{name:"show",rawName:"v-show",value:null!==t.canPR,expression:"canPR !== null"}],attrs:{id:"donate-payment"}},[t.canPR?n("div",{attrs:{id:"donate-use-pr-div"}},[n("label",{attrs:{id:"donate-use-pr-label",for:"donate-use-pr"}},[t._v("Use payment info saved "+t._s(t.deviceOrBrowser)+"?")]),n("b-form-checkbox",{attrs:{id:"donate-use-pr",switch:""},model:{value:t.usePR,callback:function(e){t.usePR=e},expression:"usePR"}})],1):t._e(),n("div",{directives:[{name:"show",rawName:"v-show",value:!t.usePR,expression:"!usePR"}]},[n("b-form-group",{staticClass:"mb-1",attrs:{label:"Your name","label-sr-only":"",state:t.nameState,"invalid-feedback":"Please enter your name."}},[n("b-form-input",{attrs:{placeholder:"Your name",autocomplete:"name",disabled:t.submitting},model:{value:t.name,callback:function(e){t.name="string"===typeof e?e.trim():e},expression:"name"}})],1),n("b-form-group",{staticClass:"mb-1",attrs:{label:"Email address","label-sr-only":"",state:t.emailState,"invalid-feedback":t.email?"This is not a valid email address.":"Please enter your email address."}},[n("b-form-input",{attrs:{type:"email",placeholder:"Email address",autocomplete:"email",disabled:t.submitting},on:{focus:function(e){t.emailFocused=!0},blur:function(e){t.emailFocused=!1}},model:{value:t.email,callback:function(e){t.email="string"===typeof e?e.trim():e},expression:"email"}})],1),n("b-form-group",{staticClass:"mb-1",attrs:{label:"Mailing address","label-sr-only":"",state:t.addressState,"invalid-feedback":"Please enter your address."}},[n("b-form-input",{attrs:{placeholder:"Mailing address",autocomplete:"street-address",disabled:t.submitting},model:{value:t.address,callback:function(e){t.address="string"===typeof e?e.trim():e},expression:"address"}})],1),n("div",{staticStyle:{display:"flex"}},[n("b-form-group",{staticClass:"mb-1",staticStyle:{flex:"auto"},attrs:{label:"City","label-sr-only":"",state:t.cityState,"invalid-feedback":"Please enter your city."}},[n("b-form-input",{attrs:{placeholder:"City",autocomplete:"address-level2",disabled:t.submitting},model:{value:t.city,callback:function(e){t.city="string"===typeof e?e.trim():e},expression:"city"}})],1),n("b-form-group",{staticClass:"mb-1",staticStyle:{flex:"none",width:"50px",margin:"0 4px"},attrs:{label:"State","label-sr-only":"",state:t.stateState,"invalid-feedback":t.state?"This is not a valid state .":"Please enter your state."}},[n("b-form-input",{attrs:{placeholder:"St",autocomplete:"address-level1",disabled:t.submitting},on:{focus:function(e){t.stateFocused=!0},blur:function(e){t.stateFocused=!1}},model:{value:t.state,callback:function(e){t.state="string"===typeof e?e.trim():e},expression:"state"}})],1),n("b-form-group",{staticClass:"mb-1",staticStyle:{flex:"none",width:"80px"},attrs:{label:"ZIP Code","label-sr-only":"",state:t.zipState,"invalid-feedback":t.zip?"This is not a valid ZIP code.":"Please enter your ZIP code."}},[n("b-form-input",{attrs:{placeholder:"ZIP",autocomplete:"postal-code",disabled:t.submitting},on:{focus:function(e){t.zipFocused=!0},blur:function(e){t.zipFocused=!1}},model:{value:t.zip,callback:function(e){t.zip="string"===typeof e?e.trim():e},expression:"zip"}})],1)],1),n("b-form-group",{staticClass:"mb-1",attrs:{label:"Card number","label-sr-only":"",state:!t.cardError&&null,"invalid-feedback":t.cardError}},[n("div",{ref:"card",staticClass:"form-control",attrs:{id:"donate-card"}})])],1),n("div",{attrs:{id:"donate-footer"}},[t.message?n("div",{attrs:{id:"donate-message"},domProps:{textContent:t._s(t.message)}}):t._e(),n("div",{attrs:{id:"donate-buttons"}},[n("b-btn",{attrs:{type:"button",variant:"secondary",disabled:t.submitting},on:{click:t.onCancel}},[t._v("Cancel")]),n("div",{directives:[{name:"show",rawName:"v-show",value:t.usePR,expression:"usePR"}],ref:"prbutton",attrs:{id:"donate-prbutton"}}),!t.usePR&&t.submitting?n("b-btn",{attrs:{id:"donate-pay-now",type:"submit",variant:"primary",disabled:""}},[n("b-spinner",{staticClass:"mr-1",attrs:{small:""}}),t._v("Paying...")],1):t._e(),t.usePR||t.submitting?t._e():n("b-btn",{attrs:{id:"donate-pay-now",type:"submit",variant:"primary"}},[t._v("Pay "+t._s(t.total?"$"+t.total/100:"Now"))])],1)])])},k=[];function P(t,e,n,a,i,r,s){try{var o=t[r](s),u=o.value}catch(c){return void n(c)}o.done?e(u):Promise.resolve(u).then(a,i)}function C(t){return function(){var e=this,n=arguments;return new Promise((function(a,i){var r=t.apply(e,n);function s(t){P(r,a,i,s,o,"next",t)}function o(t){P(r,a,i,s,o,"throw",t)}s(void 0)}))}}var _={props:{send:Function,stripeKey:String,total:Number},data:function(){return{address:"",canPR:null,card:null,cardChange:null,cardFocus:!1,city:"",elements:null,email:"",emailFocused:!1,message:null,name:"",payreq:null,prbutton:null,state:"",stateFocused:!1,submitted:!1,submitting:!1,usePR:!1,zip:"",zipFocused:!1}},mounted:function(){var t=C(d.a.mark((function t(){var e=this;return d.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return f||(f=Stripe(this.stripeKey)),this.elements=f.elements(),this.card=this.elements.create("card",{style:{base:{fontSize:"16px",lineHeight:1.5}},hidePostalCode:!0}),this.card.on("change",this.onCardChange),this.card.on("focus",(function(){e.cardFocus=!0})),this.card.on("blur",(function(){e.cardFocus=!1})),this.$nextTick((function(){e.card.mount(e.$refs.card)})),this.payreq=f.paymentRequest({country:"US",currency:"usd",total:{label:"Schola Cantorum Ticket Order",amount:100,pending:!0},requestPayerName:!0,requestPayerEmail:!0,requestShipping:!0,shippingOptions:[{id:"mail",label:"US Mail",detail:"Donation confirmation for tax records",amount:0}]}),t.next=10,this.payreq.canMakePayment();case 10:this.canPR=!!t.sent,this.canPR&&(this.usePR=!0,this.payreq.on("paymentmethod",this.onPaymentMethod),this.prbutton=this.elements.create("paymentRequestButton",{paymentRequest:this.payreq}),this.prbutton.on("click",this.onPRButtonClick),this.$nextTick((function(){e.prbutton.mount(e.$refs.prbutton)})));case 12:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),computed:{addressState:function(){return!(this.submitted&&!this.address)&&null},cardError:function(){return this.cardChange&&this.cardChange.error?this.cardChange.error.message:this.submitted?!this.cardChange||this.cardChange.empty?"Please enter your card number.":this.cardFocus||this.cardChange.complete?null:"This card number is incomplete.":null},cityState:function(){return!(this.submitted&&!this.city)&&null},deviceOrBrowser:function(){var t=navigator.userAgent||navigator.vendor;return/android|ipad|iphone|ipod|windows phone/i.test(t)?"on device":"in browser"},emailState:function(){return!(!this.emailFocused&&this.email&&!this.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/))&&(!(this.submitted&&!this.email)&&null)},nameState:function(){return(!this.submitted||""!=this.name)&&null},stateState:function(){return!(!this.stateFocused&&this.state&&!this.state.match(/^[a-zA-Z][a-zA-Z]$/))&&(!(this.submitted&&!this.state)&&null)},zipState:function(){return!(!this.zipFocused&&this.zip&&!this.zip.match(/^[0-9a-zA-Z]{5}$/))&&(!(this.submitted&&!this.zip)&&null)}},watch:{submitted:function(){this.submitted&&this.$emit("submitted")},submitting:function(){this.$emit("submitting",this.submitting)},usePR:function(){this.submitted=!1},zip:function(){this.card.update({value:{postalCode:this.zip}})}},methods:{onCancel:function(){this.$emit("cancel")},onCardChange:function(t){this.cardChange=t},onPaymentMethod:function(){var t=C(d.a.mark((function t(e){var n;return d.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return this.submitting=!0,this.card.update({disabled:!0}),t.next=4,this.send({name:e.payerName,email:e.payerEmail,address:e.shippingAddress.addressLine.join(", "),city:e.shippingAddress.city,state:e.shippingAddress.region,zip:e.shippingAddress.postalCode,subtype:"API ".concat(e.methodName),method:e.paymentMethod.id});case 4:n=t.sent,this.submitting=!1,this.card.update({disabled:!1}),n&&(this.message=n),e.complete(n?"fail":"success");case 9:case"end":return t.stop()}}),t,this)})));function e(e){return t.apply(this,arguments)}return e}(),onPRButtonClick:function(t){this.submitted=!0,this.message=null,null!==this.total?this.payreq.update({total:{label:"Schola Cantorum Ticket Order",amount:this.total,pending:!1}}):t.preventDefault()},submit:function(){var t=C(d.a.mark((function t(){var e,n,a,i;return d.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:if(null!==this.canPR&&!this.usePR){t.next=2;break}return t.abrupt("return");case 2:if(this.submitted=!0,this.message=null,null!==this.total&&null===this.nameState&&null===this.emailState&&null===this.addressState&&null===this.cityState&&null===this.stateState&&null===this.zipState&&!this.cardError){t.next=6;break}return t.abrupt("return");case 6:return this.submitting=!0,this.card.update({disabled:!0}),t.next=10,f.createPaymentMethod("card",this.card,{billing_details:{name:this.name,email:this.email,address:{line1:this.address,city:this.city,state:this.state.toUpperCase(),postal_code:this.zip}}});case 10:if(e=t.sent,n=e.paymentMethod,a=e.error,!a){t.next=19;break}return console.error(a),this.submitting=!1,this.card.update({disabled:!1}),"card_error"===a.type||"validation_error"===a.type?this.message=a.message:this.message="We’re sorry, but we're unable to process payment cards at the moment.  Please try again later, or call our office at (650) 254-1700 to order by phone.",t.abrupt("return");case 19:return t.next=21,this.send({name:this.name,email:this.email,address:this.address,city:this.city,state:this.state.toUpperCase(),zip:this.zip,subtype:"manual",method:n.id});case 21:i=t.sent,this.submitting=!1,this.card.update({disabled:!1}),i&&(this.message=i);case 25:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}()}},D=_,z=(n("c798"),Object(v["a"])(D,S,k,!1,null,null,null)),R=z.exports;function $(t,e,n,a,i,r,s){try{var o=t[r](s),u=o.value}catch(c){return void n(c)}o.done?e(u):Promise.resolve(u).then(a,i)}function O(t){return function(){var e=this,n=arguments;return new Promise((function(a,i){var r=t.apply(e,n);function s(t){$(r,a,i,s,o,"next",t)}function o(t){$(r,a,i,s,o,"throw",t)}s(void 0)}))}}var F={components:{OrderPayment:R},props:{ordersURL:String,stripeKey:String},data:function(){return{amount:0,submitted:!1,submitting:!1}},computed:{amountState:function(){return!(this.submitted&&!this.amount)&&null}},watch:{submitting:function(){this.$emit("submitting",this.submitting)}},methods:{onCancel:function(){this.$emit("cancel")},onSubmit:function(){this.$refs.pmt.submit()},onSubmitted:function(){this.submitted=!0},onSubmitting:function(t){this.submitting=t},onSend:function(){var t=O(d.a.mark((function t(e){var n,a,i,r,s,o,u,c,l;return d.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return n=e.name,a=e.email,i=e.address,r=e.city,s=e.state,o=e.zip,u=e.subtype,c=e.method,t.next=3,this.$axios.post("".concat(this.ordersURL,"/api/order"),JSON.stringify({source:"public",name:n,email:a,address:i,city:r,state:s,zip:o,lines:[{product:"donation",quantity:1,price:100*this.amount}],payments:[{type:"card",subtype:u,method:c,amount:100*this.amount}]}),{headers:{"Content-Type":"application/json"}}).catch((function(t){return t}));case 3:if(l=t.sent,!(l&&l.data&&l.data.id)){t.next=7;break}return this.$emit("success",l.data.id),t.abrupt("return",null);case 7:if(!(l&&l.data&&l.data.error)){t.next=9;break}return t.abrupt("return",l.data.error);case 9:return console.error(l),t.abrupt("return","We’re sorry, but we're unable to process payment cards at the moment.  Please try again later, or call our office at (650) 254-1700 to donate by phone.");case 11:case"end":return t.stop()}}),t,this)})));function e(e){return t.apply(this,arguments)}return e}(),setAutoFocus:function(){this.$refs.amount.focus()}}},A=F,j=(n("98b7"),Object(v["a"])(A,w,x,!1,null,null,null)),q=j.exports,T={components:{Confirmation:g,OrderForm:q},props:{ordersURL:String,products:Array,stripeKey:String,title:String},data:function(){return{orderID:null,seq:0,submitting:!1}},methods:{onClose:function(){this.$refs.modal.hide()},onHide:function(t){this.submitting&&t.preventDefault()},onOrderSuccess:function(t){this.orderID=t,this.submitting=!1},onShown:function(){this.$refs.form.setAutoFocus()},onSubmitting:function(t){this.submitting=t},show:function(){this.seq++,this.$refs.modal.show()}}},E=T,M=Object(v["a"])(E,s,o,!1,null,null,null),N=M.exports,I={props:{ordersURL:String,stripeKey:String},components:{Dialog:N}},U=I,Z=(n("9a15"),Object(v["a"])(U,i,r,!1,null,null,null)),K=Z.exports;a["default"].config.productionTip=!1,window.addEventListener("load",(function(){var t=document.getElementById("donate");t&&new a["default"]({render:function(e){return e(K,{props:{ordersURL:t.getAttribute("data-orders-url"),stripeKey:t.getAttribute("data-stripe-key")}})}}).$mount(t)}))}});
//# sourceMappingURL=donate.1de0449c.js.map