(function(t){function e(e){for(var r,i,o=e[0],l=e[1],u=e[2],d=0,p=[];d<o.length;d++)i=o[d],Object.prototype.hasOwnProperty.call(n,i)&&n[i]&&p.push(n[i][0]),n[i]=0;for(r in l)Object.prototype.hasOwnProperty.call(l,r)&&(t[r]=l[r]);c&&c(e);while(p.length)p.shift()();return s.push.apply(s,u||[]),a()}function a(){for(var t,e=0;e<s.length;e++){for(var a=s[e],r=!0,o=1;o<a.length;o++){var l=a[o];0!==n[l]&&(r=!1)}r&&(s.splice(e--,1),t=i(i.s=a[0]))}return t}var r={},n={gala:0},s=[];function i(e){if(r[e])return r[e].exports;var a=r[e]={i:e,l:!1,exports:{}};return t[e].call(a.exports,a,a.exports,i),a.l=!0,a.exports}i.m=t,i.c=r,i.d=function(t,e,a){i.o(t,e)||Object.defineProperty(t,e,{enumerable:!0,get:a})},i.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},i.t=function(t,e){if(1&e&&(t=i(t)),8&e)return t;if(4&e&&"object"===typeof t&&t&&t.__esModule)return t;var a=Object.create(null);if(i.r(a),Object.defineProperty(a,"default",{enumerable:!0,value:t}),2&e&&"string"!=typeof t)for(var r in t)i.d(a,r,function(e){return t[e]}.bind(null,r));return a},i.n=function(t){var e=t&&t.__esModule?function(){return t["default"]}:function(){return t};return i.d(e,"a",e),e},i.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},i.p="/";var o=window["webpackJsonp"]=window["webpackJsonp"]||[],l=o.push.bind(o);o.push=e,o=o.slice();for(var u=0;u<o.length;u++)e(o[u]);var c=l;s.push([2,"chunk-vendors"]),a()})({"0e6b":function(t,e,a){"use strict";var r=a("a59a"),n=a.n(r);n.a},"135e":function(t,e,a){"use strict";var r=a("2b0e"),n=a("5f5b");a("2dd8");r["default"].use(n["a"])},2:function(t,e,a){t.exports=a("e920")},9537:function(t,e,a){"use strict";var r=a("d210"),n=a.n(r);n.a},a59a:function(t,e,a){},be3b:function(t,e,a){"use strict";var r=a("2b0e"),n=a("bc3a"),s=a.n(n),i={},o=s.a.create(i);o.interceptors.request.use((function(t){return t}),(function(t){return Promise.reject(t)})),o.interceptors.response.use((function(t){return t}),(function(t){return Promise.reject(t)})),Plugin.install=function(t){t.axios=o,window.axios=o,Object.defineProperties(t.prototype,{axios:{get:function(){return o}},$axios:{get:function(){return o}}})},r["default"].use(Plugin);Plugin},d013:function(t,e,a){},d210:function(t,e,a){},e920:function(t,e,a){"use strict";a.r(e);a("744f"),a("6095"),a("6c7b"),a("d25f"),a("7514"),a("20d6"),a("f3e2"),a("1c4c"),a("6762"),a("57e7"),a("2caf"),a("cadf"),a("9865"),a("6d67"),a("e804"),a("0cd8"),a("48f8"),a("759f"),a("55dd"),a("d04f"),a("78ce"),a("8ea5"),a("0298"),a("c8ce"),a("87b3"),a("d92a"),a("217b"),a("7f7f"),a("f400"),a("7f25"),a("536b"),a("d9ab"),a("f9ab"),a("32d7"),a("25c9"),a("9f3c"),a("042e"),a("c7c6"),a("f4ff"),a("049f"),a("7872"),a("a69f"),a("0b21"),a("6c1a"),a("c7c62"),a("84b4"),a("c5f6"),a("2e37"),a("fca0"),a("7cdf"),a("ee1d"),a("b1b1"),a("87f3"),a("9278"),a("5df2"),a("04ff"),a("f751"),a("8478"),a("4504"),a("fee7"),a("1c01"),a("58b2"),a("ffc1"),a("0d6d"),a("9986"),a("8e6e"),a("25db"),a("e4f7"),a("b9a1"),a("64d5"),a("9aea"),a("db97"),a("66c8"),a("57f0"),a("165b"),a("456d"),a("cf6a"),a("fd24"),a("8615"),a("551c"),a("097d"),a("df1b"),a("2397"),a("88ca"),a("ba16"),a("d185"),a("ebde"),a("2d34"),a("f6b3"),a("2251"),a("c698"),a("a19f"),a("9253"),a("9275"),a("3b2b"),a("3846"),a("4917"),a("a481"),a("28a5"),a("386d"),a("6b54"),a("4f7f"),a("8a81"),a("ac4d"),a("8449"),a("9c86"),a("fa83"),a("48c0"),a("a032"),a("aef6"),a("d263"),a("6c37"),a("9ec8"),a("5695"),a("2fdb"),a("d0b0"),a("5df3"),a("b54a"),a("f576"),a("ed50"),a("788d"),a("14b9"),a("f386"),a("f559"),a("1448"),a("673e"),a("242a"),a("4f37"),a("c66f"),a("262f"),a("b05c"),a("34ef"),a("6aa2"),a("15ac"),a("af56"),a("b6e4"),a("9c29"),a("63d9"),a("4dda"),a("10ad"),a("c02b"),a("4795"),a("130f"),a("ac6a"),a("96cf"),a("0cdd");var r,n=a("2b0e"),s=(a("be3b"),a("135e"),function(){var t=this,e=t.$createElement,a=t._self._c||e;return t.product?a("GalaRegistration",{attrs:{product:t.product,ordersURL:t.ordersURL,stripeKey:t.stripeKey,galaRegURL:t.galaRegURL}}):t.productError?a("div",{attrs:{id:"gala-message"},domProps:{textContent:t._s(t.productError)}}):t._e()}),i=[],o=a("a34a"),l=a.n(o),u=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("b-form",{attrs:{id:"gala-form",novalidate:""},on:{submit:function(e){return e.preventDefault(),t.onSubmit(e)}}},[a("hr",{staticClass:"w-100"}),a("a",{attrs:{id:"gala-register",name:"register"}},[t._v("Registration")]),t.success?a("div",{staticClass:"mb-3"},[t._v("Thank you for your registration. A receipt has been emailed to you. We look forward to seeing you at the gala.")]):[a("b-form-group",{attrs:{label:"Register","label-for":"gala-qty","label-cols":"auto","label-class":"mt-1",state:!t.qtyError&&null,"invalid-feedback":t.qtyError}},[a("b-form-input",{staticClass:"d-inline",attrs:{id:"gala-qty",type:"number",number:"",min:"1"},model:{value:t.qty,callback:function(e){t.qty=e},expression:"qty"}}),a("span",{domProps:{textContent:t._s(t.qtyLabel)}}),a("b-form-text",[t._v("Register 8 seats to fill a table.")])],1),t._l(t.guests,(function(e,r){return a("GalaRegisterGuest",{key:r,attrs:{number:r,entreeOptions:t.entreeOptions},model:{value:t.guests[r],callback:function(e){t.$set(t.guests,r,e)},expression:"guests[i]"}})})),a("b-form-group",{attrs:{label:"Any special requests?","label-for":"gala-requests","label-class":"font-weight-bold"}},[a("b-form-textarea",{attrs:{id:"gala-requests",placeholder:"Seating preferences, dietary restrictions, etc."},model:{value:t.requests,callback:function(e){t.requests=e},expression:"requests"}})],1),a("b-form-group",{staticClass:"mb-0",attrs:{label:"Payment Information","label-class":"font-weight-bold"}}),a("OrderPayment",{ref:"pmt",attrs:{send:t.onSend,stripeKey:t.stripeKey,total:t.total,name:t.guests[0].name,email:t.guests[0].email}})]],2)},c=[],d=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("b-form-group",[a("legend",{staticClass:"bv-no-focus-ring col-form-label pt-0 font-weight-bold",attrs:{slot:"label",tabindex:"-1"},slot:"label"},[t._v("Guest "+t._s(t.number+1)),0!==t.number?a("span",{staticStyle:{"margin-left":"1em","font-weight":"normal",color:"#888","font-size":"0.875rem"}},[t._v("(leave blank if not known yet)")]):t._e()]),a("b-form-group",{attrs:{label:"Name","label-for":"guest-"+t.number+"-name","label-cols-sm":"auto","label-class":"gala-guest-label"}},[a("b-form-input",{attrs:{id:"guest-"+t.number+"-name",trim:""},model:{value:t.name,callback:function(e){t.name=e},expression:"name"}})],1),a("b-form-group",{attrs:{label:"Email","label-for":"guest-"+t.number+"-email","label-cols-sm":"auto","label-class":"gala-guest-label",state:!t.emailError&&null,"invalid-feedback":t.emailError}},[a("b-form-input",{attrs:{id:"guest-"+t.number+"-email",lazy:"",trim:"",state:!t.emailError&&null},model:{value:t.email,callback:function(e){t.email=e},expression:"email"}})],1),a("b-form-group",{attrs:{label:"Entree","label-for":"guest-"+t.number+"-entree","label-cols-sm":"auto","label-class":"gala-guest-label"}},[a("b-form-select",{attrs:{id:"guest-"+t.number+"-entree",options:t.entreeOptions},model:{value:t.entree,callback:function(e){t.entree=e},expression:"entree"}})],1)],1)},p=[],m={props:{guest:Object,number:Number,entreeOptions:Array},model:{prop:"guest",event:"input"},data:function(){return{name:"",email:"",emailError:null,entree:""}},watch:{guest:function(){this.name=this.guest.name,this.email=this.guest.email,this.entree=this.guest.entree,this.validate()},name:"emit",email:function(){this.validate(),this.emit()},entree:"emit"},methods:{emit:function(){this.$emit("input",{name:this.name,email:this.email,entree:this.entree,valid:!this.emailError})},validate:function(){this.email&&!this.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/)?this.emailError="This is not a valid email address.":this.emailError=null}}},f=m,b=(a("0e6b"),a("2877")),h=Object(b["a"])(f,d,p,!1,null,null,null),g=h.exports,y=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{directives:[{name:"show",rawName:"v-show",value:null!==t.canPR,expression:"canPR !== null"}],attrs:{id:"gala-payment"}},[t.canPR?a("div",{attrs:{id:"gala-use-pr-div"}},[a("label",{attrs:{id:"gala-use-pr-label",for:"gala-use-pr"}},[t._v("Use payment info saved "+t._s(t.deviceOrBrowser)+"?")]),a("b-form-checkbox",{attrs:{id:"gala-use-pr",switch:""},model:{value:t.usePR,callback:function(e){t.usePR=e},expression:"usePR"}})],1):t._e(),a("div",{directives:[{name:"show",rawName:"v-show",value:!t.usePR,expression:"!usePR"}]},[a("b-form-group",{staticClass:"mb-1",attrs:{label:"Your name","label-sr-only":"",state:t.nameState,"invalid-feedback":"Please enter your name."}},[a("b-form-input",{attrs:{placeholder:"Your name",autocomplete:"name",disabled:t.submitting},model:{value:t.name,callback:function(e){t.name="string"===typeof e?e.trim():e},expression:"name"}})],1),a("b-form-group",{staticClass:"mb-1",attrs:{label:"Email address","label-sr-only":"",state:t.emailState,"invalid-feedback":t.email?"This is not a valid email address.":"Please enter your email address."}},[a("b-form-input",{attrs:{type:"email",placeholder:"Email address",autocomplete:"email",disabled:t.submitting},on:{focus:function(e){t.emailFocused=!0},blur:function(e){t.emailFocused=!1}},model:{value:t.email,callback:function(e){t.email="string"===typeof e?e.trim():e},expression:"email"}})],1),a("b-form-group",{staticClass:"mb-1",attrs:{label:"Billing address","label-sr-only":"",state:t.addressState,"invalid-feedback":"Please enter your address."}},[a("b-form-input",{attrs:{placeholder:"Billing address",autocomplete:"street-address",disabled:t.submitting},model:{value:t.address,callback:function(e){t.address="string"===typeof e?e.trim():e},expression:"address"}})],1),a("div",{staticStyle:{display:"flex"}},[a("b-form-group",{staticClass:"mb-1",staticStyle:{flex:"auto"},attrs:{label:"City","label-sr-only":"",state:t.cityState,"invalid-feedback":"Please enter your city."}},[a("b-form-input",{attrs:{placeholder:"City",autocomplete:"address-level2",disabled:t.submitting},model:{value:t.city,callback:function(e){t.city="string"===typeof e?e.trim():e},expression:"city"}})],1),a("b-form-group",{staticClass:"mb-1",staticStyle:{flex:"none",width:"50px",margin:"0 4px"},attrs:{label:"State","label-sr-only":"",state:t.stateState,"invalid-feedback":t.state?"This is not a valid state .":"Please enter your state."}},[a("b-form-input",{attrs:{placeholder:"St",autocomplete:"address-level1",disabled:t.submitting},on:{focus:function(e){t.stateFocused=!0},blur:function(e){t.stateFocused=!1}},model:{value:t.state,callback:function(e){t.state="string"===typeof e?e.trim():e},expression:"state"}})],1),a("b-form-group",{staticClass:"mb-1",staticStyle:{flex:"none",width:"80px"},attrs:{label:"ZIP Code","label-sr-only":"",state:t.zipState,"invalid-feedback":t.zip?"This is not a valid ZIP code.":"Please enter your ZIP code."}},[a("b-form-input",{attrs:{placeholder:"ZIP",autocomplete:"postal-code",disabled:t.submitting},on:{focus:function(e){t.zipFocused=!0},blur:function(e){t.zipFocused=!1}},model:{value:t.zip,callback:function(e){t.zip="string"===typeof e?e.trim():e},expression:"zip"}})],1)],1),a("b-form-group",{staticClass:"mb-1",attrs:{label:"Card number","label-sr-only":"",state:!t.cardError&&null,"invalid-feedback":t.cardError}},[a("div",{ref:"card",staticClass:"form-control",attrs:{id:"gala-card"}})])],1),a("div",{attrs:{id:"gala-footer"}},[t.message?a("div",{attrs:{id:"gala-message"},domProps:{textContent:t._s(t.message)}}):t._e(),a("div",{attrs:{id:"gala-buttons"}},[a("div",{directives:[{name:"show",rawName:"v-show",value:t.usePR,expression:"usePR"}],ref:"prbutton",attrs:{id:"gala-prbutton"}}),!t.usePR&&t.submitting?a("b-btn",{attrs:{id:"gala-pay-now",type:"submit",variant:"primary",disabled:""}},[a("b-spinner",{staticClass:"mr-1",attrs:{small:""}}),t._v("Paying...")],1):t._e(),t.usePR||t.submitting?t._e():a("b-btn",{attrs:{id:"gala-pay-now",type:"submit",variant:"primary"}},[t._v("Pay "+t._s(t.total?"$"+t.total/100:"Now"))])],1)])])},v=[];function w(t,e,a,r,n,s,i){try{var o=t[s](i),l=o.value}catch(u){return void a(u)}o.done?e(l):Promise.resolve(l).then(r,n)}function x(t){return function(){var e=this,a=arguments;return new Promise((function(r,n){var s=t.apply(e,a);function i(t){w(s,r,n,i,o,"next",t)}function o(t){w(s,r,n,i,o,"throw",t)}i(void 0)}))}}var P={props:{send:Function,stripeKey:String,total:Number,name:String,email:String},data:function(){return{address:"",canPR:null,card:null,cardChange:null,cardFocus:!1,city:"",elements:null,emailFocused:!1,message:null,payreq:null,prbutton:null,state:"",stateFocused:!1,submitted:!1,submitting:!1,usePR:!1,zip:"",zipFocused:!1}},mounted:function(){var t=x(l.a.mark((function t(){var e=this;return l.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return r||(r=Stripe(this.stripeKey)),this.elements=r.elements(),this.card=this.elements.create("card",{style:{base:{fontSize:"16px",lineHeight:1.5}},hidePostalCode:!0}),this.card.on("change",this.onCardChange),this.card.on("focus",(function(){e.cardFocus=!0})),this.card.on("blur",(function(){e.cardFocus=!1})),this.$nextTick((function(){e.card.mount(e.$refs.card)})),this.payreq=r.paymentRequest({country:"US",currency:"usd",total:{label:"Schola Cantorum Gala Registration",amount:100,pending:!0},requestPayerName:!0,requestPayerEmail:!0,requestShipping:!0,shippingOptions:[{id:"mail",label:"US Mail",detail:"Donation confirmation for tax records",amount:0}]}),t.next=10,this.payreq.canMakePayment();case 10:this.canPR=!!t.sent,this.canPR&&(this.usePR=!0,this.payreq.on("paymentmethod",this.onPaymentMethod),this.prbutton=this.elements.create("paymentRequestButton",{paymentRequest:this.payreq}),this.prbutton.on("click",this.onPRButtonClick),this.$nextTick((function(){e.prbutton.mount(e.$refs.prbutton)})));case 12:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),computed:{addressState:function(){return!(this.submitted&&!this.address)&&null},cardError:function(){return this.cardChange&&this.cardChange.error?this.cardChange.error.message:this.submitted?!this.cardChange||this.cardChange.empty?"Please enter your card number.":this.cardFocus||this.cardChange.complete?null:"The card information is incomplete.":null},cityState:function(){return!(this.submitted&&!this.city)&&null},deviceOrBrowser:function(){var t=navigator.userAgent||navigator.vendor;return/android|ipad|iphone|ipod|windows phone/i.test(t)?"on device":"in browser"},emailState:function(){return!(!this.emailFocused&&this.email&&!this.email.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/))&&(!(this.submitted&&!this.email)&&null)},nameState:function(){return(!this.submitted||""!=this.name)&&null},stateState:function(){return!(!this.stateFocused&&this.state&&!this.state.match(/^[a-zA-Z][a-zA-Z]$/))&&(!(this.submitted&&!this.state)&&null)},zipState:function(){return!(!this.zipFocused&&this.zip&&!this.zip.match(/^[0-9a-zA-Z]{5}$/))&&(!(this.submitted&&!this.zip)&&null)}},watch:{submitted:function(){this.submitted&&this.$emit("submitted")},submitting:function(){this.$emit("submitting",this.submitting)},usePR:function(){this.submitted=!1},zip:function(){this.card.update({value:{postalCode:this.zip}})}},methods:{onCardChange:function(t){this.cardChange=t},onPaymentMethod:function(){var t=x(l.a.mark((function t(e){var a;return l.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return this.submitting=!0,this.card.update({disabled:!0}),t.next=4,this.send({name:e.payerName,email:e.payerEmail,address:e.shippingAddress.addressLine.join(", "),city:e.shippingAddress.city,state:e.shippingAddress.region,zip:e.shippingAddress.postalCode,subtype:"API ".concat(e.methodName),method:e.paymentMethod.id});case 4:a=t.sent,this.submitting=!1,this.card.update({disabled:!1}),a&&(this.message=a),e.complete(a?"fail":"success");case 9:case"end":return t.stop()}}),t,this)})));function e(e){return t.apply(this,arguments)}return e}(),onPRButtonClick:function(t){this.submitted=!0,this.message=null,null!==this.total?this.payreq.update({total:{label:"Schola Cantorum Gala Registration",amount:this.total,pending:!1}}):t.preventDefault()},submit:function(){var t=x(l.a.mark((function t(){var e,a,n,s;return l.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:if(null!==this.canPR&&!this.usePR){t.next=2;break}return t.abrupt("return");case 2:if(this.submitted=!0,this.message=null,null!==this.total&&null===this.nameState&&null===this.emailState&&null===this.addressState&&null===this.cityState&&null===this.stateState&&null===this.zipState&&!this.cardError){t.next=6;break}return t.abrupt("return");case 6:return this.submitting=!0,this.card.update({disabled:!0}),t.next=10,r.createPaymentMethod("card",this.card,{billing_details:{name:this.name,email:this.email,address:{line1:this.address,city:this.city,state:this.state.toUpperCase(),postal_code:this.zip}}});case 10:if(e=t.sent,a=e.paymentMethod,n=e.error,!n){t.next=19;break}return console.error(n),this.submitting=!1,this.card.update({disabled:!1}),"card_error"===n.type||"validation_error"===n.type?this.message=n.message:this.message="We’re sorry, but we're unable to process payment cards at the moment.  Please try again later, or call our office at (650) 254-1700 to order by phone.",t.abrupt("return");case 19:return t.next=21,this.send({name:this.name,email:this.email,address:this.address,city:this.city,state:this.state.toUpperCase(),zip:this.zip,subtype:"manual",method:a.id});case 21:s=t.sent,this.submitting=!1,this.card.update({disabled:!1}),s&&(this.message=s);case 25:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}()}},S=P,R=(a("f4e0"),Object(b["a"])(S,y,v,!1,null,null,null)),k=R.exports;function q(t,e,a,r,n,s,i){try{var o=t[s](i),l=o.value}catch(u){return void a(u)}o.done?e(l):Promise.resolve(l).then(r,n)}function z(t){return function(){var e=this,a=arguments;return new Promise((function(r,n){var s=t.apply(e,a);function i(t){q(s,r,n,i,o,"next",t)}function o(t){q(s,r,n,i,o,"throw",t)}i(void 0)}))}}function C(t,e){return A(t)||E(t,e)||_()}function _(){throw new TypeError("Invalid attempt to destructure non-iterable instance")}function E(t,e){if(Symbol.iterator in Object(t)||"[object Arguments]"===Object.prototype.toString.call(t)){var a=[],r=!0,n=!1,s=void 0;try{for(var i,o=t[Symbol.iterator]();!(r=(i=o.next()).done);r=!0)if(a.push(i.value),e&&a.length===e)break}catch(l){n=!0,s=l}finally{try{r||null==o["return"]||o["return"]()}finally{if(n)throw s}}return a}}function A(t){if(Array.isArray(t))return t}var O={components:{GalaRegisterGuest:g,OrderPayment:k},props:{ordersURL:String,galaRegURL:String,product:Object,stripeKey:String},data:function(){return{qty:1,qtyError:null,qtyValid:!0,guests:[{name:"",email:"",entree:"",valid:!0}],requests:"",success:!1}},computed:{entreeOptions:function(){var t=[{text:"(please select)",value:""}];return this.product.options.forEach((function(e){var a=e.split("/",2),r=C(a,2),n=r[0],s=r[1];t.push({text:n,value:s})})),t},qtyLabel:function(){return"number"===typeof this.qty&&1===this.qty?" seat at $".concat(this.product.price/100):" seats at $".concat(this.product.price/100," each")},total:function(){return this.valid?this.qty*this.product.price:null},valid:function(){return!this.qtyError&&!this.guests.find((function(t){return!t.valid}))}},watch:{qty:function(){if(this.qtyValid="number"===typeof this.qty&&this.qty>=1,this.qtyValid)if(this.qtyError=null,this.qty<this.guests.length)this.guests.splice(this.qty);else for(var t=this.guests.length;t<this.qty;t++)this.guests.push({name:"",email:"",entree:"",valid:!0});else this.qtyError="Please enter a valid quantity."}},methods:{onSend:function(){var t=z(l.a.mark((function t(e){var a,r,n,s,i,o,u,c,d,p,m=this;return l.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return a=e.name,r=e.email,n=e.address,s=e.city,i=e.state,o=e.zip,u=e.subtype,c=e.method,d=new FormData,d.append("source","public"),d.append("name",a),d.append("email",r),d.append("address",n),d.append("city",s),d.append("state",i),d.append("zip",o),d.append("cNote",this.requests),this.guests.forEach((function(t,e){var a="line".concat(e+1,".");d.append(a+"product",m.product.id),d.append(a+"quantity",1),d.append(a+"price",m.product.price),d.append(a+"guestName",t.name),d.append(a+"guestEmail",t.email),d.append(a+"option",t.entree)})),d.append("payment1.type","card"),d.append("payment1.subtype",u),d.append("payment1.method",c),d.append("payment1.amount",this.total),t.next=17,this.$axios.post(this.galaRegURL,d,{headers:{"Content-Type":"application/x-www-form-urlencoded"}}).catch((function(t){return t}));case 17:if(p=t.sent,!(p&&p.data&&p.data.id)){t.next=21;break}return this.success=!0,t.abrupt("return",null);case 21:if(!(p&&p.data&&p.data.error)){t.next=23;break}return t.abrupt("return",p.data.error);case 23:return console.error(p),t.abrupt("return","We’re sorry, but we're unable to process payment cards at the moment.  Please try again later, or call our office at (650) 254-1700 to register by phone.");case 25:case"end":return t.stop()}}),t,this)})));function e(e){return t.apply(this,arguments)}return e}(),onSubmit:function(){this.$refs.pmt.submit()}}},$=O,j=(a("9537"),Object(b["a"])($,u,c,!1,null,null,null)),Z=j.exports;function F(t,e,a,r,n,s,i){try{var o=t[s](i),l=o.value}catch(u){return void a(u)}o.done?e(l):Promise.resolve(l).then(r,n)}function U(t){return function(){var e=this,a=arguments;return new Promise((function(r,n){var s=t.apply(e,a);function i(t){F(s,r,n,i,o,"next",t)}function o(t){F(s,r,n,i,o,"throw",t)}i(void 0)}))}}var L={components:{GalaRegistration:Z},props:{ordersURL:String,galaRegURL:String,stripeKey:String,productID:String},data:function(){return{product:null,productError:null}},created:function(){var t=U(l.a.mark((function t(){var e;return l.a.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,this.$axios.get("".concat(this.ordersURL,"/payapi/prices?p=").concat(this.productID));case 2:if(e=t.sent.data,e){t.next=5;break}return t.abrupt("return");case 5:"string"===typeof e?this.productError=e:this.product=e.products[0];case 6:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}()},T=L,N=Object(b["a"])(T,s,i,!1,null,null,null),I=N.exports;n["default"].config.productionTip=!1,window.addEventListener("load",(function(){var t=document.getElementById("gala-registration");if(t){var e=t.getAttribute("data-product");new n["default"]({render:function(a){return a(I,{props:{ordersURL:t.getAttribute("data-orders-url"),galaRegURL:t.getAttribute("data-gala-reg-url"),stripeKey:t.getAttribute("data-stripe-key"),productID:e}})}}).$mount(t)}}))},f4e0:function(t,e,a){"use strict";var r=a("d013"),n=a.n(r);n.a}});
//# sourceMappingURL=gala.c1019ec8.js.map