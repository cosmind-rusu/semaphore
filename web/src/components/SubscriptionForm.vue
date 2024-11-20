<template>
  <v-form
    ref="form"
    lazy-validation
    v-model="formValid"
    v-if="item != null"
  >
    <v-alert
      :value="formError"
      :type="(formError || '').includes('already activated') ? 'warning' : 'error'"
    >{{ formError }}
    </v-alert>

    <v-text-field
      v-model="item.key"
      label="Subscription Key"
      :rules="[v => !!v || $t('key_required')]"
      required
      :disabled="formSaving"
    ></v-text-field>

    <v-list v-if="item.plan">
      <v-list-item class="pa-0">
        <v-list-item-content>
          <v-list-item-title>Plan</v-list-item-title>
          <v-list-item-subtitle>{{ item.plan }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-list-item class="pa-0">
        <v-list-item-content>
          <v-list-item-title>Status</v-list-item-title>
          <v-list-item-subtitle>
            <v-chip :color="statusColor" label class="mt-1">{{ item.state }}</v-chip>
          </v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-list-item class="pa-0">
        <v-list-item-content>
          <v-list-item-title>Expires at</v-list-item-title>
          <v-list-item-subtitle>{{ item.expiresAt }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-list-item class="pa-0">
        <v-list-item-content>
          <v-list-item-title>Users</v-list-item-title>
          <v-list-item-subtitle>{{ item.users }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
    </v-list>
    <v-alert type="info" v-else>There is no active subscription.</v-alert>

  </v-form>
</template>
<script>
import ItemFormBase from '@/components/ItemFormBase';

export default {
  mixins: [ItemFormBase],

  data() {
    return {};
  },

  computed: {
    isNew() {
      return false;
    },

    statusColor() {
      switch (this.item.state) {
        case 'expired':
          return 'error';
        case 'active':
          return 'success';
        default:
          return '';
      }
    },
  },

  methods: {
    async afterSave() {
      await this.loadData();
    },

    getItemsUrl() {
      return '/api/subscription';
    },

    getSingleItemUrl() {
      return '/api/subscription';
    },

    getRequestOptions() {
      return {
        method: 'post',
        url: '/api/subscription',
      };
    },
  },
};
</script>
