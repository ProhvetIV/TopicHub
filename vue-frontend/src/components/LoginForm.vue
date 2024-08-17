<template>
    <div id="pageContent">
        <div class="loginWrapper">
            <h1>Login</h1>
            <form @submit.prevent="submitForm">
                <div class="login-input">
                    <label for="username">Username:</label>
                    <input v-model="formData.username" type="text" id="username" required />
                </div>
                <div class="login-input">
                    <label for="password">Password:</label>
                    <input v-model="formData.password" type="password" id="password" required />
                </div>
                <p :class="{ 'invisible': !badLoginVisible, 'margins': badLoginVisible }">{{ badLoginText }}</p>
                <button type="submit" class="nice-button">Submit</button>
            </form>
        </div>
    </div>
</template>
  
<script setup>
    import { ref, watch } from 'vue';
    import { checkIfFilled, sendMessageWhenOpen } from '@/utils.js'; // @ = shorthand alias for the src directory
    import { EventBus } from '@/eventBus';
    
    const formData = ref({
        username: '',
        password: ''
    });
    const badLoginVisible = ref(false);
    const badLoginText = 'Wrong username or password';


    const submitForm = () => {
        if (!checkIfFilled(formData.value)) {
            badLoginVisible.value = false;
            sendMessageWhenOpen('login', formData.value);

            //const webSocketService = getWebSocketService();            
            //webSocketService.sendMessage('login', formData.value);
        } else {
            badLoginVisible.value = true;
        }
    };

    watch(() => EventBus.badLogin, () =>{
        badLoginVisible.value = true;
    })
</script>
  
<style scoped>
    form {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: space-around;
        row-gap: 10px;
    }
    .login-input {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;
    }
    .invisible {
        display: none;
    }
    .margins {
        margin: 10px 0;
    }
</style>