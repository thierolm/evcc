<template>
	<div>
		<button
			v-show="iconVisible"
			href="#"
			data-bs-toggle="modal"
			data-bs-target="#notificationModal"
			class="btn btn-link text-decoration-none link-light text-nowrap"
		>
			<fa-icon :class="iconClass" icon="exclamation-triangle"></fa-icon>
		</button>

		<div
			id="notificationModal"
			class="modal fade"
			tabindex="-1"
			role="dialog"
			aria-hidden="true"
		>
			<div
				class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable"
				role="document"
			>
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">{{ $t("notifications.modalTitle") }}</h5>
						<button
							type="button"
							class="btn-close"
							data-bs-dismiss="modal"
							aria-label="Close"
						></button>
					</div>
					<div class="modal-body">
						<div v-for="(msg, index) in notifications" :key="index">
							<small
								class="d-flex justify-content-end mt-3"
								:title="fmtAbsoluteDate(msg.time)"
							>
								{{ fmtTimeAgo(msg.time) }}
							</small>
							<p class="d-flex align-items-baseline">
								<fa-icon
									:class="{
										'text-danger': msg.type === 'error',
										'text-warning': msg.type === 'warn',
									}"
									class="flex-grow-0 d-block"
									icon="exclamation-triangle"
								></fa-icon>
								<span class="flex-grow-1 px-2 py-1 text-break">
									{{ msg.message }}
								</span>
								<span class="badge rounded-pill bg-secondary" v-if="msg.count > 1">
									{{ msg.count }}
								</span>
							</p>
						</div>
					</div>
					<div class="modal-footer">
						<button
							type="button"
							data-bs-dismiss="modal"
							aria-label="Close"
							@click="clear"
							class="btn btn-outline-secondary"
						>
							{{ $t("notifications.dismissAll") }}
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import "../icons";
import formatter from "../mixins/formatter";

export default {
	name: "Notifications",
	props: {
		notifications: Array,
	},
	computed: {
		iconVisible: function () {
			return this.notifications.length > 0;
		},
		iconClass: function () {
			return this.notifications.find((m) => m.type === "error")
				? "text-danger"
				: "text-warning";
		},
	},
	methods: {
		clear: function () {
			window.app && window.app.clear();
		},
	},
	created: function () {
		this.interval = setInterval(() => {
			this.$forceUpdate();
		}, 10 * 1000);
	},
	destroyed: function () {
		clearTimeout(this.interval);
	},
	mixins: [formatter],
};
</script>
