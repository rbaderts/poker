
'use strict';


var card =  Vue.component('card', {
	name: "card",
	template: `
		<img :alt="imageSrc" :src="imageSrc" 
		v-bind:class="{'pcard': typeA, 'pcard1': typeB, 'pcard2': typeC, 'ocard1': typeD, 'ocard2': typeE }"/> `,
	props: {
		index: {
			type: Number,
			default: 0
		},
		suit: {
			type: Number,
			default: 1
		},
		typeA: {
			type: Boolean,
			default: false,
		},
		typeB: {
			type: Boolean,
			default: false,
		},
		typeC: {
			type: Boolean,
			default: false,
		},
		typeD: {
			type: Boolean,
			default: false,
		},
		typeE: {
			type: Boolean,
			default: false,
		}
	},
	computed: {
		imageSrc() {
			if ((this.index  == -1) || (this.suit == -1)) {
				return "/static/img/BackRed_Folded.png";
			} else if ((this.index  == 0) || (this.suit == 0)) {
				return "/static/img/BackRed.png";
			}
			return '/static/img/Classic/'+this.cardSuit(this.suit) + this.cardSymbol(this.index) + ".png"
		}
	} ,
	methods : {
		cardSuit: function (suit) {
			if (suit == 1) {
				return 's';
			} else if (suit == 2) {
				return 'h';
			} else if (suit == 3) {
				return 'd';
			} else {
				return 'c';
			}
		},
		cardSymbol: function(index) {
			let str1 = Number(index).toString();
			let str = str1.padStart(2, '0');
			return str;
		}
	}
});
/*
var card =  Vue.component('card', {
	name: "card",
	template: `
		<div class="card" :class="[suitSym, 'values-' + indexSym]" 
		style="background-image: "image">
			<h1>{{indexSym}}</h1>
			<div class="figures" :class="suitSym"></div>
			<h1>{{indexSym}}</h1>
		</div> `,
	props: ['index', 'suit'],
	computed: {
		indexSym: function () {
			return this.cardSymbol(this.index);
		},
		suitSym: function () {
			return this.cardSuit(this.suit);
		},

	},
	methods: {
		cardSuit: function(suit) {
			if (suit == 1) {
				return 'figures-P';
			} else if (suit == 2) {
				return 'figures-H';
			} else if (suit == 3) {
				return 'figures-D';
			} else if (suit == 4) {
				return 'figures-C';
			}
		},
		cardSymbol: function(index) {
			if (index < 10) {
				return index
			} else if (index == 11) {
				return 'J';
			} else if (index == 12) {
				return 'J';
			} else if (index == 13) {
				return 'K';
			} else if (index == 14) {
				return 'A';
			}
		}


	}
});
*/
