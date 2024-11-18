<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div>
    <v-dialog
      v-model="paymentDialog"
      max-width="400"
      persistent
      :transition="false"
    >
      <v-card
        :disabled="paypalButtonRendering"
        v-if="payment == null || payment.gateway === 'paypal' && payment.state !== 'completed'"
      >
        <v-card-title class="headline text-center">
          Replenishing your wallet
          <v-spacer></v-spacer>
          <v-btn
            @click="closePaymentDialog"
            icon
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>
        <v-card-text class="text-xs-center pb-0">

          <v-slider
            v-model="currencyAmount"
            always-dirty
            :min="projectType === 'premium' ? 25 : 5"
            :max="projectType === 'premium' ? 250 : 50"
            :step="projectType === 'premium' ? 25 : 5"
            thumb-label="always"
            thumb-size="60"
            style="margin-top: 90px; margin-left: 15px; margin-right: 15px;"
          >
            <template v-slot:thumb-label="props">
              <div
                style="font-size: 20px; font-weight: bold;"
              >
                ${{ props.value }}
              </div>
            </template>
          </v-slider>

        </v-card-text>
        <v-card-actions class="pb-4 pt-0 d-block">

          <div class="mt-4" id="paypal-button-container"></div>

        </v-card-actions>
      </v-card>

      <v-card v-else>
        <v-card-title class="headline text-center">
          Adding funds
          <v-spacer></v-spacer>
          <v-btn
            @click="closePaymentDialog"
            icon
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>

        <v-alert
          :value="paymentError"
          color="error"
        >
          {{ paymentError }}
        </v-alert>

        <v-card-text class="text-center">
          <div v-if="payment.state === 'completed'">
            <v-icon class="ma-2" color="success" style="font-size: 86px;">mdi-check-circle</v-icon>
            <div class="title pb-1">Payment completed</div>
            <div class="">Thank you for your payment.</div>
          </div>

          <div v-else>
            <v-progress-circular
              color="primary"
              :size="70"
              :width="7"
              indeterminate
              class="mb-3 mt-3"
            ></v-progress-circular>
            <div class="title pb-1">Awaiting payment</div>
            <div class="">Complete payment in the pop-up window.</div>
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>

    <YesNoDialog
      title="New subscription key"
      text="Are you sure?"
      v-model="regenerateKeyDialog"
      @yes="regenerateKey()"
    />

    <v-dialog
      v-model="activationGuideDialog"
      max-width="400"
      persistent
      :transition="false"
      scrollable
    >
      <v-card>
        <v-card-title>
          How to activate subscription?
          <v-spacer></v-spacer>
          <v-btn
            @click="activationGuideDialog = false"
            icon
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>

        <v-divider></v-divider>

        <v-card-text class="text-xs-center pt-3" style="max-height: 1000px;">
          <p class="mb-2">
            1. Download and install <b>Semaphore Pro</b>.
          </p>
          <v-btn
            color="primary"
            class="mb-4"
            :to="`/project/${projectId}/install`"
          >Download and Install</v-btn>

          <p class="mb-2">
            2. Log in as an administrator to the Semaphore UI and click the yellow
            button labeled <b>Activate Pro Subscription</b>.
          </p>
          <v-img class="rounded mb-6" :aspect-ratio="489/351"
                 src="activation_guide/screen1.webp?v=2"></v-img>

          <p class="mb-2">
            3. Enter your subscription key
            <code>{{ (project || {}).licenseKey }} <v-btn
              v-if="(project || {}).licenseKey"
              class="ml-1"
              icon
              x-small
              @click="copyToClipboard((project || {}).licenseKey)"
            >
              <v-icon>mdi-content-copy</v-icon>
            </v-btn>
            </code>
            in the dialog box that appears and click the <b>Activate</b> button.
          </p>
          <v-img class="rounded mb-6" :aspect-ratio="428/277"
                 src="activation_guide/screen2.webp"></v-img>

          <p class="mb-2">
            4. After successful activation, your subscription details will be displayed.
          </p>
          <v-img class="rounded" :aspect-ratio="428/408"
                 src="activation_guide/screen3.webp"></v-img>
        </v-card-text>
      </v-card>
    </v-dialog>

    <v-toolbar flat>
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>
        {{ $t('dashboard') }}
        <SubscriptionLabel v-if="projectType === 'premium'" />
      </v-toolbar-title>
    </v-toolbar>

    <DashboardMenu
      :project-id="projectId"
      :project-type="projectType"
      :can-update-project="true"
    />

    <v-container v-if="project != null" class="ml-0 px-6">
      <div class="mt-7 mb-6">
        <div class="text-h3">
          <v-icon
            x-large
            color="deep-orange darken-4"
            style="margin-top: -8px;"
          >
            mdi-wallet
          </v-icon>

          ${{ project.balance }}

          <v-btn
            color="primary"
            @click="paymentDialog = true"
            icon
            x-large
            style="margin-top: -8px;"
          >
            <v-icon x-large>mdi-plus-circle</v-icon>
          </v-btn>
        </div>
      </div>

      <div>
        <v-timeline
          align-top
          dense
          style="margin-left: -20px;"
        >
          <v-timeline-item
            fill-dot
            icon="mdi-calendar-range"
            class="text-subtitle-1 align-center"
          >
            Billing date:
            <span v-if="project.planCanceled">&mdash;</span>
            <span v-else>{{ project.planFinishDate | formatDate2 }}</span>
          </v-timeline-item>

          <v-timeline-item
            v-if="projectType === 'premium'"
            fill-dot
            icon="mdi-professional-hexagon"
            class="text-subtitle-1 align-center"
          >
            Subscription Key: <code>{{ project.licenseKey || '&mdash;' }}</code>
            <v-btn
              v-if="project.licenseKey"
              class="ml-2"
              icon
              @click="copyToClipboard(project.licenseKey)"
            >
              <v-icon>mdi-content-copy</v-icon>
            </v-btn>

            <v-btn
              v-if="project.licenseKey"
              class="ml-2"
              icon
              @click="regenerateKeyDialog = true"
            >
              <v-icon>mdi-refresh</v-icon>
            </v-btn>

            <a
              v-if="project.licenseKey"
              class="ml-4"
              @click="activationGuideDialog = true"
            >How to use it?</a>

          </v-timeline-item>

          <v-timeline-item
            v-if="projectType === 'premium' && isPremiumActive"
            fill-dot
            icon="mdi-download"
            class="text-subtitle-1 align-center"
          >
            <v-btn
              color="primary"
              :to="`/project/${projectId}/install`"
            >Download
            </v-btn>

            <a class="ml-4" @click="activationGuideDialog = true">How to activate it?</a>
          </v-timeline-item>

          <v-timeline-item
            v-if="projectType === ''"
            fill-dot
            icon="mdi-server"
            class="text-subtitle-1 align-center"
          >Cache: {{ project.diskUsage }} / {{ plan.diskUsage }} Mb used
          </v-timeline-item>
          <v-timeline-item
            v-if="projectType === ''"
            fill-dot
            icon="mdi-cog"
            class="text-subtitle-1 align-center"
          >Runner:
            {{ Math.ceil(project.runnerUsage / 60) }} / {{ Math.ceil(plan.runnerUsage / 60) }}
            minutes used
          </v-timeline-item>
        </v-timeline>
      </div>

      <v-row class="mt-0 mb-9" v-if="projectType === ''">
        <v-col md="4" lg="4">
          <v-card
            class="mt-4 pa-2"
            :color="$vuetify.theme.dark ? 'blue-grey darken-4' : 'grey lighten-4'"
            flat
          >
            <v-card-title class="text-h3">
              Free
              <v-spacer/>
              <v-chip
                class="ml-4 mt-1 text-subtitle-1 py-3 px-4 font-weight-bold"
                v-if="project.plan === 'free'"
                color="success"
                outlined
              >Your plan
              </v-chip>

            </v-card-title>
            <v-card-text>
              <v-timeline
                align-top
                dense
                style="margin-left: -30px;"
              >
                <v-timeline-item
                  icon="mdi-server"
                  class="text-subtitle-1 align-center"
                  fill-dot
                >
                  100 Mb for cache
                </v-timeline-item>

                <v-timeline-item
                  fill-dot
                  icon="mdi-cog"
                  class="text-subtitle-1 align-center"
                >
                  50 min/mo for tasks
                </v-timeline-item>
              </v-timeline>
            </v-card-text>
            <v-card-actions>
              <v-btn large style="visibility: hidden;">
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
        <v-col md="4" lg="4">
          <v-card
            class="mt-4 pa-2"
            :color="$vuetify.theme.dark ? 'blue-grey darken-4' : 'grey lighten-4'"
            flat
          >
            <v-card-title class="text-h3">
              $7

              <v-spacer/>
              <v-chip
                class="ml-4 mt-1 font-weight-bold text-center pa-4 text-subtitle-1"
                color="success"
                outlined
                :style="{'font-size': project.planCanceled ? '12px !important' : undefined}"
                v-if="project.plan === 'starter'"
              >
                Your plan
                <span class="ml-1" v-if="project.planCanceled">
                  until {{ project.planFinishDate | formatDate2 }}
                </span>
              </v-chip>
            </v-card-title>
            <v-card-text>

              <v-timeline
                align-top
                dense
                style="margin-left: -30px;"
              >
                <v-timeline-item
                  icon="mdi-server"
                  fill-dot
                  class="text-subtitle-1 align-center"
                >
                  1G for cache
                </v-timeline-item>

                <v-timeline-item
                  fill-dot
                  icon="mdi-cog"
                  class="text-subtitle-1 align-center"
                >
                  <div>1000 min/mo for tasks</div>
                </v-timeline-item>
              </v-timeline>
            </v-card-text>

            <v-card-actions>
              <v-btn
                depressed
                large
                :color="getPlanColor(project.plan)"
                style="width: 100%;"
                @click="
                project.plan === 'free' || project.planCanceled
                  ? selectPlan('starter')
                  : selectPlan('free')
              "
              >
                {{
                  project.plan === 'free'
                    ? 'Upgrade'
                    : (project.planCanceled ? 'Continue Auto Renew' : 'Cancel Auto Renew')
                }}
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
      <v-row class="mt-0 mb-9" v-else-if="projectType === 'premium'">
        <v-col v-for="plan in premiumPlans" md="4" lg="4" :key="plan.id">
          <v-card
            class="mt-4 pa-2"
            :color="$vuetify.theme.dark ? 'blue-grey darken-4' : 'grey lighten-4'"
            flat
          >
            <v-card-title class="text-h3">
              ${{ plan.price }}
              <v-spacer/>
              <v-chip
                class="text-subtitle-1 py-3 px-4 font-weight-bold"
                v-if="project.plan === plan.id"
                :style="{'font-size': project.planCanceled ? '12px !important' : undefined}"
                color="success"
                outlined
              >Your plan
                <span class="ml-1" v-if="project.planCanceled">
                  until {{ project.planFinishDate | formatDate2 }}
                </span>
              </v-chip>
            </v-card-title>

            <v-card-text>
              <v-timeline
                align-top
                dense
                style="margin-left: -30px;"
              >
                <v-timeline-item
                  fill-dot
                  icon="mdi-account-multiple"
                  class="text-subtitle-1 align-center"
                >
                  <div>{{ plan.users }} users</div>
                </v-timeline-item>
              </v-timeline>
            </v-card-text>

            <v-card-actions>
              <v-btn
                depressed
                large
                :color="getPlanColor(plan.id)"
                style="width: 100%;"
                @click="
                project.plan !== plan.id || project.planCanceled
                  ? selectPlan(plan.id)
                  : selectPlan('free')
              "
              >
                {{ getActivePremiumButtonTitle(plan.id) }}
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>

    </v-container>
  </div>
</template>
<style lang="scss">
</style>
<script>
import EventBus from '@/event-bus';
import axios from 'axios';
import { loadScript } from '@paypal/paypal-js';
import { getErrorMessage } from '@/lib/error';
import YesNoDialog from '@/components/YesNoDialog.vue';
import SubscriptionLabel from '@/components/SubscriptionLabel.vue';
import DashboardMenu from '@/components/DashboardMenu.vue';

const PLANS = {
  free: {
    price: 0,
    diskUsage: 100,
    runnerUsage: 50 * 60,
  },
  starter: {
    price: 7,
    diskUsage: 1000,
    runnerUsage: 1000 * 60,
  },
  pro: {
    price: 25,
    diskUsage: 3000,
    runnerUsage: 3000 * 60,
  },
  premium: {
    price: 12,
    users: 4,
  },
  premium_plus: {
    price: 40,
    users: 8,
  },
  pro1: {
    price: 5,
    users: 1,
  },
  pro4: {
    price: 15,
    users: 4,
  },
  pro10: {
    price: 30,
    users: 10,
  },
  enterprise: {
    price: 1000,
    users: 1000000,
  },
};

const FREE_TRIAL_PLANS = ['home', 'premium'];

export default {
  components: { YesNoDialog, SubscriptionLabel, DashboardMenu },
  props: {
    projectId: Number,
    projectType: String,
  },

  data() {
    return {
      project: null,
      payment: null,
      paymentError: null,
      paymentDialog: false,
      paymentProgressTimer: null,
      currencyAmount: null,
      plan: PLANS.free,
      /**
       * @array
       */
      premiumPlans: ['pro1', 'pro4', 'pro10'].map((plan) => ({
        ...PLANS[plan],
        id: plan,
      })),
      paypal: null,
      paypalButton: null,
      paypalButtonRendering: true,
      regenerateKeyDialog: null,
      activationGuideDialog: null,
    };
  },

  computed: {
    isPremiumActive() {
      return this.premiumPlans.findIndex((p) => p.id === this.project.plan) >= 0;
    },
  },

  watch: {
    projectId() {
      this.refreshProject();
    },

    async paymentDialog(val) {
      if (val) {
        this.currencyAmount = null;
        await this.initPaypalButton();
      } else {
        this.paypalButton.close();
      }
    },
  },

  async created() {
    await this.refreshProject();
  },

  methods: {
    async regenerateKey() {
      const { data } = (await axios({
        method: 'post',
        url: `/billing/projects/${this.projectId}/subscription/rekey`,
        responseType: 'json',
        data: {},
      }));

      this.project.licenseKey = data.key;
    },

    getPlanColor(plan) {
      if (plan === 'free') {
        return 'success';
      }

      if (plan !== this.project.plan) {
        return 'success';
      }

      if (this.project.planCanceled) {
        return 'warning';
      }

      return 'error';
      // return project.plan === 'free' || project.planCanceled ? 'success' : 'error';
    },

    async initPaypalButton() {
      if (!this.paypal) {
        try {
          this.paypal = await loadScript({ clientId: 'ATj8K7c7xzD4Z1JMozXCQGh6sUv8M9yScCeQWfVm1xaC_8UQ_0AM7IU6kr1cFbQ3FKwq0_dOz4-gkE_k' });
        } catch (error) {
          console.error('failed to load the PayPal JS SDK script', error);
        }
      }

      try {
        this.paypalButtonRendering = true;
        this.paypalButton = await this.paypal.Buttons({
          createOrder: async () => {
            try {
              this.payment = (await axios({
                method: 'post',
                url: `/billing/projects/${this.projectId}/payments`,
                responseType: 'json',
                data: {
                  currencyAmount: this.currencyAmount,
                  currency: 'usd',
                  gateway: 'paypal',
                },
              })).data;

              console.log(this.payment);

              return this.payment.gatewayTransactionId;
            } catch (err) {
              EventBus.$emit('i-snackbar', {
                color: 'error',
                text: getErrorMessage(err),
              });
              return undefined;
            }
          },

          onApprove: async () => {
            this.payment = (await axios({
              method: 'put',
              url: `/billing/projects/${this.projectId}/payments/${this.payment.number}/refresh`,
              responseType: 'json',

              data: {
                currencyAmount: this.currencyAmount,
                currency: 'usd',
                gateway: 'paypal',
              },
            })).data;

            await this.refreshProject();
          },
        });

        await this.paypalButton.render('#paypal-button-container');
      } catch (error) {
        console.error('failed to render the PayPal Buttons', error);
      } finally {
        this.paypalButtonRendering = false;
      }
    },

    async copyToClipboard(text) {
      try {
        await window.navigator.clipboard.writeText(text);
        EventBus.$emit('i-snackbar', {
          color: 'success',
          text: 'The Subscription Key has been copied to the clipboard.',
        });
      } catch (e) {
        EventBus.$emit('i-snackbar', {
          color: 'error',
          text: `Can't copy the Subscription Key: ${e.message}`,
        });
      }
    },

    getActivePremiumButtonTitle(planId) {
      if (this.project.plan === 'free') {
        return this.project.freeTrialAvailable && FREE_TRIAL_PLANS.includes(planId)
          ? 'Try free trial'
          : 'Activate';
      }

      if (this.project.plan === planId) {
        if (this.project.planCanceled) {
          return 'Continue Auto Renew';
        }
        return 'Cancel Auto Renew';
      }

      if (PLANS[this.project.plan].price < PLANS[planId].price) {
        return 'Upgrade';
      }

      return 'Downgrade';
    },

    showDrawer() {
      EventBus.$emit('i-show-drawer');
    },

    onError(e) {
      EventBus.$emit('i-snackbar', {
        color: 'error',
        text: e.message,
      });
    },

    async refreshProject() {
      this.project = (await axios({
        method: 'get',
        url: `/billing/projects/${this.projectId}`,
        responseType: 'json',
      })).data;

      this.plan = PLANS[this.project.plan];

      if (this.plan && this.premiumPlans.some((plan) => plan.id === this.project.plan)) {
        this.premiumPlans.push({
          ...this.plan,
          id: this.project.plan,
        });
      }
    },

    async selectPlan(plan) {
      await this.refreshProject();

      const freeTrial = this.project.freeTrialAvailable && FREE_TRIAL_PLANS.includes(plan);

      const { price } = PLANS[plan];

      if (!freeTrial && this.project.balance < price) {
        this.paymentDialog = true;
      } else {
        await axios({
          method: 'put',
          url: `/billing/projects/${this.projectId}`,
          responseType: 'json',
          data: {
            plan,
          },
        });
        await this.refreshProject();
      }
    },

    async refreshPayment() {
      this.payment = (await axios({
        method: 'get',
        url: `/billing/projects/${this.projectId}/payments/${this.payment.number}`,
        responseType: 'json',
      })).data;
    },

    async closePaymentDialog() {
      this.payment = null;
      this.paymentError = null;
      this.paymentDialog = false;
    },

    async makeCoinbasePayment() {
      window.gtag_report_conversion(this.currencyAmount);

      try {
        this.payment = (await axios({
          method: 'post',
          url: `/billing/projects/${this.projectId}/payments`,
          responseType: 'json',
          // headers: {
          //   authorization: `Bearer ${localStorage.getItem('authenticationToken')}`,
          // },
          data: {
            currencyAmount: this.currencyAmount,
            currency: 'usd',
            gateway: 'coinbase',
          },
        })).data;

        this.paymentError = null;
        this.paymentDialog = true;

        // eslint-disable-next-line no-promise-executor-return
        await new Promise((resolve) => setTimeout(resolve, 600));

        window.open(this.payment.hostedUrl, '_blank');

        this.paymentProgressTimer = setInterval(async () => {
          await this.refreshPayment();
          if (this.payment.state !== 'new' && this.payment.state !== 'pending') {
            clearInterval(this.paymentProgressTimer);
            await this.refreshProject();
          }
        }, 2000);
      } catch (err) {
        EventBus.$emit('i-snackbar', {
          color: 'error',
          text: getErrorMessage(err),
        });
      }
    },
  },
};
</script>
