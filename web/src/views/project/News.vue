<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
    <div>
        <v-toolbar flat>
            <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
            <v-toolbar-title>
                {{ $t('dashboard') }}
                <SubscriptionLabel v-if="projectType === 'premium'" />
            </v-toolbar-title>
        </v-toolbar>
    </div>
</template>
<script>
export default {};
</script>
