<template>
  <v-tabs show-arrows class="pl-4">

    <v-tab
      key="artifacts"
      :to="`/project/${projectId}/install`"
    >
      Artifacts
    </v-tab>

  </v-tabs>
</template>
<script>

export default {

  props: {
    projectId: Number,
  },

  data() {
    return {
      id: null,
    };
  },

};
</script>
