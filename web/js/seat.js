'use strict';

//import card from './card.js'

var seat =  Vue.component('seat', {
	name: "seat",
	props:{
		seatnum:{
			type:Number,
			default:0
		},
		stack:{
			type:Number,
			default:0
		},
		ontable:{
			type:Number,
			default:0
		},
		username: {
			type: String,
			default: ""
		},
		ourSeatNumber: {
			type: Number,
			default: 0
		},
		color: {
			type: String,
			default: ""
		},
		card1: {
			type: Object,
			default: null
		},
		card2: {
			type: Object,
			default: null
		},
		folded: {
			type: Boolean,
			default: false 
		},
		dealer: {
			type: Boolean,
			default: false
		},
		percentwin: {
			type: Number,
			default: 0
		},
		percenttie: {
			type: Number,
			default: 0
		},
		hand: {
			type: String,
			default: ""
		}
	},
	data() {
			return {
				message: "",
				startTime: 0,
				interval: null,
				sendmessage: ""
			}
	},
	model:{
	},
		/*
	created: function () {
		this.attach()
	},
	*/
    mounted: function () {
	    new CircleType(document.getElementById('player-name'))
		    .radius(48);
    },

	methods: {

		sendMessage: function() {
			let vm = this;
			let msg = vm.sendmessage;
			let tableid = window.tableId;
			axios.post(`/api/tables/`+tableid+`/seat/`+vm.seatnum+'/message', msg);
		},

		newHand: function() {
			this.message = "";
			this.card1 = null;
			this.card2 = null;
			this.ontable = 0;
		},
		fold: function() {
		},
		chipCount: function (n, next, ontable) {
			let x = ontable;
			if (next != 0) {
				x = x % next;
			}
			let f = Math.round((x / n) + 0.5) - 1;
			let result = f.toFixed(0)
			console.log("chipCount n=" + n + "next=," + next + ", ontable="+ontable + ",result="+result);
			return result;
		},
		showPlayerMessage: function (message) {
			this.message = message;
			this.startMessageTimer();
		},
		startMessageTimer: function() {
			var vm = this;

			let now = new Date();
			vm.startTime = now.getTime();

			console.log("start message timer");
			this.interval = setInterval(function() {
				let now = new Date();
				let elapsed = now.getTime() - vm.startTime;
				if (elapsed > 10000) {
					vm.message = "";
					console.log("end message timer");
					clearInterval(vm.interval)
				}
			}, 1000);

		},
	},


	computed: {
		//ontable: function() {
		//	return this.initialontable + this.addontable
		//},
		playerPosition: function () {
			//return ["player-" + (this.player.seatnum + 1), {"playing": this.player.seatnum === this.ourseatnumber}]
			let offset = this.ourSeatNumber;

//			let playing = "";
//			return  "player-" + (this.seatnum+1) + " " + playing;
			let playing = "playing";
			if (this.seatnum === this.ourSeatNumber) {
				//return  "player-" + (this.seatnum+1) + " " + playing;
				let playing = "playing";
				return  "player-1" + " " + playing;
			}
			//return  "player-" + (this.seatnum+1) + " " + playing;
			return  "player-" + (this.seatnum+1) + " " + playing;
		}
	},
	template: `
	<div class="player">
	
	<!--

// <h2 id="demo1">Here’s some curved text flowing clockwise.</h2>
new CircleType(document.getElementById('demo1'))
  .radius(384);
-->

		<div class="avatar" :style="{backgroundColor: color || 'dodgerblue'}">
		     <div id="player-name" class="player-name">{{username}}</div>
			 <div class="player-seatnum">{{ seatnum }}</div>
    		 <div class="player-bank-value">{{ stack }}</div>
     		 <img class="player-dealerbutton" v-if="dealer == true" src="/static/img/dealerbutton.png"/>

		 </div>
        <div class="player-cards">
			<card v-if="card1 != null" :index="card1.index" :suit="card1.suit" :type-b="true"/>
			<card v-else-if="folded == true" :index="-1" :suit="-1" :type-b="true"/>
			<card v-else :index="0" :suit="0" :type-b="true"/>
			<card v-if="card2 != null" :index="card2.index" :suit="card2.suit" :type-c="true"/>
			<card v-else-if="folded == true" :index="-1" :suit="-1" :type-c="true"/>
			<card v-else :index="0" :suit="0"  :type-c="true"/>
		</div>
		
		<div v-if="ontable > 0 && folded == false" class="player-bet">
			<chipstack  v-bind:ontable="ontable"/>
		</div>
		
		<div class="hand-description">
		    <i>{{ hand }}</i>
		</div>
		
		<!--
		<div v-if="folded == true">
      		<h3>folded</h3>
		</div>
		-->
		
		<div class="odds-card">
		     <strong class="odds-name">Win</strong> <strong class="odds-value"> {{ percentwin }}% </strong> <br>
		     <strong class="odds-name">Tie</strong> <strong class="odds-value"> {{ percenttie }}% </strong>
		</div>
		
	</div>`
	});

