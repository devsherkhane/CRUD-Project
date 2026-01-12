<template>
  <div class="form-session">
    <h2>Session 1: Student Registration</h2>
    <form @submit.prevent="submitForm">
      <input v-model="student.studentName" placeholder="Name" required />
      <input v-model="student.email" type="email" placeholder="Email" required />
      <input v-model="student.mobileNumber" placeholder="Mobile" required />
      <input v-model="student.address" placeholder="Address" />

      <select v-model="student.state" @change="onStateChange" required>
        <option value="">Select State</option>
        <option v-for="(districts, stateName) in locations" :key="stateName" :value="stateName">
          {{ stateName }}
        </option>
      </select>

      <select v-model="student.district" @change="onDistrictChange" :disabled="!student.state" required>
        <option value="">Select District</option>
        <option v-for="(talukas, distName) in availableDistricts" :key="distName" :value="distName">
          {{ distName }}
        </option>
      </select>

      <select v-model="student.taluka" :disabled="!student.district" required>
        <option value="">Select Taluka</option>
        <option v-for="taluka in availableTalukas" :key="taluka" :value="taluka">
          {{ taluka }}
        </option>
      </select>

      <select v-model="student.gender">
        <option value="">Select Gender</option>
        <option value="Male">Male</option>
        <option value="Female">Female</option>
      </select>
      
      <input type="date" v-model="student.dob" />
      <input v-model="student.bloodGroup" placeholder="Blood Group" />
      
      <div class="file-input">
        <label>Student Photo:</label>
        <input type="file" @change="handleFileUpload" accept="image/*" />
      </div>

      <label>
        <input type="checkbox" v-model="student.handicapped" /> Handicapped
      </label>
      
      <button type="submit">Submit Data</button>
    </form>
    <router-link to="/list">View Registered Students</router-link>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { locations } from '../data/locations'; // Import the local data

const student = ref({
  studentName: '', email: '', mobileNumber: '', address: '',
  state: '', district: '', taluka: '', gender: '',
  dob: '', handicapped: false, bloodGroup: ''
});

const photoFile = ref(null);

// Computed: Get districts for the selected state
const availableDistricts = computed(() => {
  return student.value.state ? locations[student.value.state] : {};
});

// Computed: Get talukas for the selected district
const availableTalukas = computed(() => {
  if (student.value.state && student.value.district) {
    return locations[student.value.state][student.value.district] || [];
  }
  return [];
});

// Resets child dropdowns when the parent changes
const onStateChange = () => {
  student.value.district = '';
  student.value.taluka = '';
};

const onDistrictChange = () => {
  student.value.taluka = '';
};

const handleFileUpload = (event) => {
  photoFile.value = event.target.files[0];
};

const submitForm = async () => {
  try {
    const formData = new FormData();
    // Use student object directly for the body
    Object.keys(student.value).forEach(key => {
      formData.append(key, student.value[key]);
    });
    
    if (photoFile.value) {
      formData.append('photo', photoFile.value);
    }

    const response = await fetch('http://localhost:8000/students', {
      method: 'POST',
      body: formData
    });

    if (response.ok) {
      alert('Student added successfully!');
      // Reset form instead of reloading page for better UX
      student.value = {
        studentName: '', email: '', mobileNumber: '', address: '',
        state: '', district: '', taluka: '', gender: '',
        dob: '', handicapped: false, bloodGroup: ''
      };
      photoFile.value = null;
    }
  } catch (err) {
    console.error("Error submitting form:", err);
  }
};
</script>
<style>
  .form-session {
  max-width: 500px;
  margin: 0 auto;
  background: white;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

form {
  display: flex;
  flex-direction: column;
  gap: 15px; /* Adds space between every input */
}

input, select {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

input:focus, select:focus {
  outline: none;
  border-color: #42b983;
  box-shadow: 0 0 5px rgba(66, 185, 131, 0.3);
}

.file-input {
  background: #f9f9f9;
  padding: 10px;
  border: 1px dashed #ccc;
  border-radius: 4px;
}

button[type="submit"] {
  background-color: #42b983;
  color: white;
  padding: 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background 0.3s ease;
}

button[type="submit"]:hover {
  background-color: #3aa876;
}
</style>