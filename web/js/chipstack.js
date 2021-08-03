
'use strict';


var chipstack =  Vue.component('chipstack', {
	name: "chipstack",
	props: {
		ontable: {
			type: Number,
			default: 0
		}
	},
	data() {
		return {
			notcomplete: true
		}
	},

template: `<div>

       <div class="ontable-value">{{ ontable }}</div> 
       <div class="ontable">
				<div class="stack-10000">
					<div class="stack v-10000"
						  v-for="(n, i) in ((ontable - (ontable % 10000)) / 10000)"
						 :style="{top: ((-2 + (i * 5))) + 'px'}"
						 v-if="ontable >= 10000"/>
				</div>
				<div class="stack-1000">
					<div class="stack v-1000"
						  v-for="(n, i) in ((ontable%10000 - (ontable%10000)%1000) / 1000)"
						 :style="{top: ((-2 + (i * 5))) + 'px'}"
						 v-if="ontable % 10000 >= 1000"/>
				</div>
				<div class="stack-100">
					<div class="stack v-100"
						  v-for="(n, i) in ((ontable%1000 - (ontable%1000)%100) / 100)"
						 :style="{top: ((-2 + (i * 5))) + 'px'}"
						 v-if="ontable % 1000 >= 100"/>
				</div>
				<div class="stack-25">
					<div class="stack v-25"
						  v-for="(n, i) in (((ontable%100) - ((ontable%100)%25)) / 25)"
						 :style="{top: ((-2 + (i * 5))) + 'px'}"
						 v-if="ontable % 100 >= 25"/>
				</div>
				<div class="stack-5">
					<div class="stack v-5" 
						v-for="(n, i) in (((ontable%25) - ((ontable%25)%5)) / 5)"
						 :style="{top: ((-2 + (i * 5))) + 'px'}"
						v-if="ontable % 25 >= 5"/>
				</div>
				<div class="stack-1">
					<div class="stack v-1" 
						v-for="(n, i) in (ontable%5)"
						 :style="{top: ((-2 + (i * 5))) + 'px'}"
						v-if="ontable % 5 > 0"/>
				</div>
<!--				<div class="ontable-value">{{ ontable }}</div>-->

                </div>

      	     </div>`,
	methods : {
	}
});

