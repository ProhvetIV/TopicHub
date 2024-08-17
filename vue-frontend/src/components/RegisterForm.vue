<template>
    <div id="pageContent">
        <div class="registerWrapper">
            <h1>Register</h1>
            <form id="registration-form" @submit.prevent="submitForm">
                <div>
                    <label for="nickname">Nickname:</label>
                    <input v-model="formData.nickname" type="text" id="nickname">
                </div>
                <div>
                    <label for="ageEl">Age:</label>
                    <input v-model="formData.age" type="number" id="ageEl" required>
                </div>
                <div class="div-radio-buttons">
                    <label>Gender:</label>
                    <div v-for="gender in genders" :key="gender">
                        <label :for="gender">{{ gender }}</label>
                        <input v-model="formData.gender" type="radio" :value="gender" :id="gender">
                    </div>
                </div>
                <div>
                    <label for="firstName">First Name:</label>
                    <input v-model="formData.firstName" type="text" id="firstName" required>
                </div>
                <div>
                    <label for="lastName">Last Name:</label>
                    <input v-model="formData.lastName" type="text" id="lastName" required>
                </div>
                <div>
                    <label for="email">E-mail:</label>
                    <input v-model="formData.email" type="text" id="email" required>
                </div>
                <div>
                    <label for="username">Username:</label>
                    <input v-model="formData.username" type="text" id="username">
                </div>
                <div>
                    <label for="password">Password:</label>
                    <input v-model="formData.password" type="password" id="password" required>
                </div>
                <div>
                    <label for="aboutMe">About me:</label>
                    <textarea class="aboutMe"  v-model="formData.aboutMe" type="text"></textarea>
                </div>

                <div v-if="imageSrc !== null" class="flex-div">
                    <img class="profile-image" :src="imageSrc">
                </div>
                
                <div class="flex-div">
                    <button type="button" class="userProfileButtons" @click="triggerFileInput">Upload File</button>
                    <input type="file" ref="fileInput" @change="handleFileUpload" accept=".jpeg, .jpg, .png, .gif, .webp" style="display: none;"/>
                </div>

                <p :class="{ 'invisible': !badLoginVisible, 'margins': badLoginVisible }">{{ badLoginText }}</p>

                <div class="flex-div">
                    <button type="submit" class="nice-button register-submit">Submit</button>                
                </div>
            </form>
        </div>
    </div>
</template>
  
<script setup>
    import { ref } from 'vue';
    import { checkIfFilled, sendMessageWhenOpen } from '@/utils.js'; // @ = shorthand alias for the src directory
    import { getWebSocketService } from '@/websocket.js';

    let imageSrc = ref(null);
    const fileInput = ref(null);
    
    const formData = ref({
        nickname: '',
        age: '',
        gender: '',
        firstName: '',
        lastName: '',
        email: '',
        username: '',
        password: '',
        imageData: '',
        imageName: '',
        aboutMe: '',
    })

    const triggerFileInput = () => {
        console.log(fileInput.value); // Debugging line
        if (fileInput.value) {
            fileInput.value.click();
        } else {
            console.error("fileInput is not defined or click is not a function");
        }
    };

    const handleFileUpload = (event) => {
        const file = event.target.files[0];
        if (file) {
            imageSrc.value = URL.createObjectURL(file);

            // Perform actions with the file, e.g., upload it to a server or display it
            /*console.log('Selected file:', file);
            console.log('Selected file parts:', file.name, file.webkitRelativePath);
            getWebSocketService().sendMessage("postProfilePic", {image: file});*/
            const reader = new FileReader();
            reader.onload = () => {
                const base64String = reader.result.split(',')[1]; // Get Base64 string without prefix
                console.log("file.name: ", file.name);
                
                formData.value.imageData = base64String;
                formData.value.imageName = file.name;         
            };
            reader.readAsDataURL(file);
        }
    };



    const genders = ['male', 'female', 'other']
    const badLoginVisible = ref(false);
    const badLoginText = ref('Please fill in all fields.');
    const submitForm = () => {
        //if (!checkIfFilled(formData.value)) {
        if (!sortOptionalFields(formData.value)) {
            if (formData.value.imageData === '') {
                formData.value.imageData = null
            }
            badLoginVisible.value = false;
            sendMessageWhenOpen('UserRegistration', formData.value);
            //getWebSocketService().sendMessage('UserRegistration', formData.value);
        } else {
            badLoginVisible.value = true;
        }
    }

    const sortOptionalFields = (formData) => {        
        const { nickname, username, imageData, imageName, aboutMe, ...requiredData } = formData;
        console.log("formdata: ", formData, "requireddata: ", requiredData);
//
        // Send the requiredData object to the validation function
        return checkIfFilled(requiredData.value)
    }
</script>

<style scoped>
    #registration-form div {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }

    .div-radio-buttons {
        flex-direction: column;
    }

    .div-radio-buttons div {
        width: 30%;
    }

    .flex-div {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        margin: 10px 0px;
    }

    .register-submit {
        margin-left: 10px;
    }

    .aboutMe{
        border: 1px solid #ccc;
        border-radius: 5px;
        width: 192px;
        height: 100px;
        resize: none;
    }
</style>