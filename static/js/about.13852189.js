(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["about"],{"90ca":function(t,e,i){"use strict";i.r(e);var s=function(){var t=this,e=t.$createElement,i=t._self._c||e;return i("section",{staticClass:"section section-picture-details"},[t.pictureDetails?i("div",[i("div",{staticClass:"box"},[i("figure",{staticClass:"image is-16by9"},[i("img",{attrs:{alt:t.pictureDetails.title,src:t.getAbsoluteURL(t.pictureDetails.image.highres_url)}}),i("div",{staticClass:"tags has-addons image-tags"},[i("span",{staticClass:"tag is-dark is-medium"},[i("b-icon",{attrs:{icon:"comments"}})],1),i("span",{staticClass:"tag is-info is-medium"},[t._v(t._s(t.pictureDetails.num_comments))]),i("span",{staticClass:"tag is-dark is-medium"},[i("b-icon",{attrs:{icon:"thumbs-up"}})],1),i("span",{staticClass:"tag is-info is-medium"},[t._v(t._s(t.pictureDetails.num_likes))])]),i("div",{staticClass:"buttons image-buttons"},[t.userLikesPic?i("b-button",{staticClass:"is-uppercase",attrs:{type:"is-danger","icon-left":"times-circle"},on:{click:t.unLikePicture}},[t._v("\n            Stop Liking\n          ")]):i("b-button",{staticClass:"is-uppercase",attrs:{type:"is-success","icon-left":"thumbs-up"},on:{click:t.likePicture}},[t._v("\n            Like Picture\n          ")])],1)]),i("br"),i("h1",{staticClass:"title"},[t._v(t._s(t.pictureDetails.title))]),i("p",{staticClass:"is-size-4"},[t._v(t._s(t.pictureDetails.description))]),i("hr",{staticClass:"hr"}),i("h2",{staticClass:"title"},[t._v("User Comments")]),i("b-message",{attrs:{type:"is-danger","has-icon":"",size:"small",active:t.pictureComments&&0===t.pictureComments.length}},[t._v("\n        There are no comments for this picture yet\n      ")]),t.pictureComments&&t.pictureComments.length>0?i("div",{staticClass:"comments-wrapper"},t._l(t.pictureComments,(function(e){return i("div",{key:"comment_id_"+e.comment_id,staticClass:"media"},[i("div",{staticClass:"media-left"},[e.user&&e.user.avatar?i("figure",{staticClass:"image avatar is-64x64"},[i("img",{attrs:{alt:e.user.fullname,src:t.getAbsoluteURL(e.user.avatar.thumb_url)}})]):t._e()]),i("div",{staticClass:"media-content is-size-5"},[t._v(t._s(e.comment))])])})),0):t._e(),i("form-wrapper",{attrs:{validator:t.$v.newComment}},[i("form",{on:{submit:function(e){return e.preventDefault(),t.createComment(e)}}},[i("form-group",{attrs:{name:"comment",label:"Post new comment"}},[i("b-input",{attrs:{type:"textarea",placeholder:"You comment for this picture..."},on:{input:function(e){return t.$v.newComment.comment.$touch()}},model:{value:t.newComment.comment,callback:function(e){t.$set(t.newComment,"comment",e)},expression:"newComment.comment"}})],1),i("div",{staticClass:"field"},[i("b-button",{attrs:{type:"is-primary",size:"is-large","icon-left":"comment"},on:{click:t.postComment}},[t._v("\n              Post Comment\n            ")])],1),i("b-message",{attrs:{type:"is-danger","has-icon":"",size:"small",active:t.commentError},on:{"update:active":function(e){t.commentError=e}}},[t._v("\n            Could not create the comment\n          ")])],1)])],1)]):t._e(),i("b-loading",{attrs:{"is-full-page":!0,active:t.isLoading},on:{"update:active":function(e){t.isLoading=e}}})],1)},n=[],r=(i("96cf"),i("3b8d")),a=i("b5ae"),u={name:"PictureDetails",data:function(){return{isLoading:!0,pictureDetails:null,pictureComments:null,pictureLikes:null,userLikesPic:!1,userID:0,commentError:!1,newComment:{picture_id:null,user_id:null,comment:null}}},created:function(){this.userID=parseInt(this.$ls.get("user_id")),this.loadPictureDetails(),this.$eventHub.$on("picture-created",this.goToPictures)},methods:{goToPictures:function(){this.$router.push({name:"pictures"})},loadPictureDetails:function(){var t=Object(r["a"])(regeneratorRuntime.mark((function t(){var e;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return this.isLoading=!0,t.next=3,this.apiGet("pictures/"+this.$route.params.id);case 3:if(e=t.sent,200!==e.status){t.next=12;break}return this.pictureDetails=e.data,t.next=8,this.loadPictureComments();case 8:return t.next=10,this.loadPictureLikes();case 10:t.next=13;break;case 12:this.goToPictures();case 13:this.isLoading=!1;case 14:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),loadPictureComments:function(){var t=Object(r["a"])(regeneratorRuntime.mark((function t(){var e;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,this.apiGet("pictures/"+this.$route.params.id+"/comments");case 2:e=t.sent,200===e.status?this.pictureComments=e.data:this.pictureComments=[];case 4:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),loadPictureLikes:function(){var t=Object(r["a"])(regeneratorRuntime.mark((function t(){var e,i=this;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,this.apiGet("pictures/"+this.$route.params.id+"/likes");case 2:e=t.sent,200===e.status?this.pictureLikes=e.data:this.pictureLikes=[],this.userLikesPic=this.pictureLikes.filter((function(t){return t.user_id===i.userID})).length>0;case 5:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),unLikePicture:function(){var t=Object(r["a"])(regeneratorRuntime.mark((function t(){var e;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,this.apiDelete("pictures/"+this.$route.params.id+"/likes/"+this.userID);case 2:e=t.sent,204===e.status?this.loadPictureDetails():this.$buefy.toast.open({duration:5e3,message:"Error updating Picture likes",position:"is-bottom",type:"is-danger"});case 4:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),likePicture:function(){var t=Object(r["a"])(regeneratorRuntime.mark((function t(){var e;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,this.apiPost("pictures/"+this.$route.params.id+"/likes/"+this.userID);case 2:e=t.sent,201===e.status?(this.loadPictureDetails(),this.$buefy.snackbar.open("Picture Liked!")):this.$buefy.toast.open({duration:5e3,message:"Error updating Picture likes",position:"is-bottom",type:"is-danger"});case 4:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),postComment:function(){var t=Object(r["a"])(regeneratorRuntime.mark((function t(){var e;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:if(this.newComment.user_id=this.userID,this.newComment.picture_id=parseInt(this.$route.params.id),this.$v.newComment.$touch(),!1!==this.$v.newComment.$invalid){t.next=8;break}return t.next=6,this.apiPost("pictures/"+this.$route.params.id+"/comments",this.newComment);case 6:e=t.sent,201===e.status?(this.loadPictureDetails(),this.$buefy.snackbar.open("You comment has been created successfully"),this.newComment.comment=null,this.$v.newComment.$reset()):this.$buefy.toast.open({duration:5e3,message:"Error creating new comment",position:"is-bottom",type:"is-danger"});case 8:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}()},validations:{newComment:{comment:{required:a["required"],minLength:Object(a["minLength"])(5),maxLength:Object(a["maxLength"])(2e3)}}}},c=u,o=i("2877"),m=Object(o["a"])(c,s,n,!1,null,null,null);e["default"]=m.exports},e1f1:function(t,e,i){"use strict";i.r(e);var s=function(){var t=this,e=t.$createElement,i=t._self._c||e;return i("section",{staticClass:"section section-pictures"},[t.pictures&&0===t.pictures.length?i("b-message",{attrs:{title:"No Pictures",type:"is-info","has-icon":""}},[t._v("\n    There are no Pictures created yet\n  ")]):t._e(),i("div",{staticClass:"columns is-multiline is-centered"},t._l(t.pictures,(function(e){return i("div",{key:"pic_"+e.picture_id,staticClass:"column is-4-mobile is-3-desktop"},[i("router-link",{attrs:{to:{name:"picture",params:{id:e.picture_id}}}},[i("figure",{staticClass:"image is-16by9",staticStyle:{background:"#999999"}},[i("img",{attrs:{alt:e.title,src:t.getAbsoluteURL(e.image.lowres_url)}})])])],1)})),0),i("b-loading",{attrs:{"is-full-page":!0,active:t.isLoading},on:{"update:active":function(e){t.isLoading=e}}})],1)},n=[],r=(i("96cf"),i("3b8d")),a={name:"Pictures",data:function(){return{isLoading:!0,pictures:null}},created:function(){this.loadPictures(),this.$eventHub.$on("picture-created",this.onPictureCreated)},methods:{loadPictures:function(){var t=Object(r["a"])(regeneratorRuntime.mark((function t(){var e;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return this.isLoading=!0,t.next=3,this.apiGet("pictures");case 3:e=t.sent,200===e.status?this.pictures=e.data:204===e.status?this.pictures=[]:this.$router.push({name:"home"}),this.isLoading=!1;case 6:case"end":return t.stop()}}),t,this)})));function e(){return t.apply(this,arguments)}return e}(),onPictureCreated:function(){this.loadPictures()}}},u=a,c=i("2877"),o=Object(c["a"])(u,s,n,!1,null,null,null);e["default"]=o.exports}}]);
//# sourceMappingURL=about.13852189.js.map