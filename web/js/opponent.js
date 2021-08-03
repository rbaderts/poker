'use strict';

//import card from './card.js'

var opponent =  Vue.component('opponent', {
	name: "opponent",
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
		ourSeatNumber: {
			type: Number,
			default: 0
		},
		username: {
			type: String,
			default: ""
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
		numseats: {
			type: Number,
			default: 0
		},
     	dealer: {
			type: Boolean,
			default: false
		},
		folded: {
			type: Boolean,
			default: false
		}
	},
	data() {
			return {
				callAmount: 0,
				requestId: 0,
				stackAvailable: 0,
				currentBet: 0,
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
	mounted: function () {
		let name = document.getElementById(this.playerPosition);
		if (name != null) {
			new CircleType(name).radius(40);
		}
	},

	computed: {
		playerPosition: function () {

			//  0 1 2(o) 3 4  :   opponents:    1 2 3 4

			// len = 5:      len - (seatnum - oruseatNum)

			/// if s > our.
//			1 2 1 3    ==>

			let pos = 0;
			if (this.seatnum > this.ourSeatNumber) {
				pos = this.numseats - (this.seatnum - this.ourSeatNumber) + 1;
			} else if (this.seatnum < this.ourSeatNumber) {
				pos = (this.numseats - this.ourSeatNumber - 1 - this.seatnum) + 1;
			}

			if (pos <= 0) {
				pos = this.numseats + pos;

			}

    		return  "opponent-" + pos
		}
	},
	template: `
	<div class="opponent" :class="playerPosition">

		<div class="avatar" :style="{backgroundColor: color || 'dodgerblue'}">
		     <div :id="playerPosition" class="opponent-player-name">{{username}} </div>
		     <div class="opponent-seatnum">{{ seatnum }}</div>
		     <div class="opponent-bank-value">{{ stack }}</div>
  	       	 <img class="opponent-dealerbutton" v-if="dealer == true" src="/static/img/dealerbutton.png"/>
        </div>
        
        <div class="opponent-cards">
			<card v-if="folded == true" :index="-1" :suit="-1" :type-d="true"/>
			<card v-else-if="card1 != null" :index="card1.index" :suit="card1.suit" :type-d="true"/>
			<card v-else :index="0" :suit="0" :type-d="true"/>
			<card v-if="folded == true" :index="-1" :suit="-1" :type-e="true"/>
			<card v-else-if="card2 != null" :index="card2.index" :suit="card2.suit" :type-e="true"/>
			<card v-else :index="0" :suit="0" :type-e="true"/>
		</div>
		
		<div v-if="ontable > 0" class="opponent-bet">
			<chipstack  v-bind:ontable="ontable"/>
		</div>
		
		
		

	</div>`
	});

