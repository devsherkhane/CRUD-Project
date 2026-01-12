<template>
  <div class="list-session">
    <div class="header-section">
      <h2>Session 2: Student Records</h2>
      <router-link to="/" class="nav-link">Add New Student</router-link>
    </div>

    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Photo</th>
            <th>Name</th>
            <th>Email</th>
            <th>Mobile</th>
            <th>Address</th>
            <th>State</th>
            <th>District</th>
            <th>Taluka</th>
            <th>Gender</th>
            <th>DOB</th>
            <th>Blood</th>
            <th>H.C.</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in students" :key="s.id">
            <td class="photo-cell">
              <img 
                v-if="s.photo" 
                :src="'http://localhost:8000/uploads/' + s.photo" 
                alt="Student"
                class="student-img"
              />
              <div v-else class="no-photo">N/A</div>
            </td>
            
            <td>
              <input v-if="editId === s.id" v-model="editData.studentName" />
              <span v-else>{{ s.studentName }}</span>
            </td>

            <td>
              <input v-if="editId === s.id" v-model="editData.email" type="email" />
              <span v-else>{{ s.email }}</span>
            </td>

            <td>
              <input v-if="editId === s.id" v-model="editData.mobileNumber" />
              <span v-else>{{ s.mobileNumber }}</span>
            </td>

            <td>
              <input v-if="editId === s.id" v-model="editData.address" />
              <span v-else>{{ s.address }}</span>
            </td>

            <td>
              <input v-if="editId === s.id" v-model="editData.state" />
              <span v-else>{{ s.state }}</span>
            </td>

            <td>
              <input v-if="editId === s.id" v-model="editData.district" />
              <span v-else>{{ s.district }}</span>
            </td>

            <td>
              <input v-if="editId === s.id" v-model="editData.taluka" />
              <span v-else>{{ s.taluka }}</span>
            </td>

            <td>
              <select v-if="editId === s.id" v-model="editData.gender">
                <option value="Male">Male</option>
                <option value="Female">Female</option>
              </select>
              <span v-else>{{ s.gender }}</span>
            </td>

            <td>
              <input v-if="editId === s.id" v-model="editData.dob" type="date" />
              <span v-else>{{ s.dob }}</span>
            </td>

            <td>
              <input v-if="editId === s.id" v-model="editData.bloodGroup" class="small-input" />
              <span v-else>{{ s.bloodGroup }}</span>
            </td>

            <td class="center-text">
              <input v-if="editId === s.id" v-model="editData.handicapped" type="checkbox" />
              <span v-else>{{ s.handicapped ? 'Yes' : 'No' }}</span>
            </td>

            <td>
              <div v-if="editId === s.id" class="action-buttons">
                <button class="save-btn" @click="updateStudent(s.id)">Save</button>
                <button class="cancel-btn" @click="editId = null">Cancel</button>
              </div>
              <div v-else class="action-buttons">
                <button class="edit-btn" @click="startEdit(s)">Edit</button>
                <button class="delete-btn" @click="deleteStudent(s.id)">Delete</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const students = ref([]);
const editId = ref(null);
const editData = ref({});

const fetchStudents = async () => {
  const res = await fetch('http://localhost:8000/students');
  students.value = await res.json();
};

const deleteStudent = async (id) => {
  if (confirm('Are you sure you want to delete this record?')) {
    await fetch(`http://localhost:8000/students/${id}`, { method: 'DELETE' });
    fetchStudents();
  }
};

const startEdit = (student) => {
  editId.value = student.id;
  editData.value = { ...student };
};

const updateStudent = async (id) => {
  await fetch(`http://localhost:8000/students/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(editData.value)
  });
  editId.value = null;
  fetchStudents();
};

onMounted(fetchStudents);
</script>

<style scoped>
.list-session {
  max-width: 1300px;
  margin: 20px auto;
  padding: 20px;
  font-family: 'Segoe UI', sans-serif;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

h2 { color: #2c3e50; margin: 0; }

.nav-link {
  text-decoration: none;
  background: #42b983;
  color: white;
  padding: 8px 16px;
  border-radius: 4px;
  font-weight: 600;
}

.table-container {
  overflow-x: auto;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

table {
  width: 100%;
  border-collapse: collapse;
  min-width: 1000px;
}

th {
  background-color: #f8f9fa;
  color: #555;
  text-align: left;
  padding: 15px;
  font-size: 13px;
  border-bottom: 2px solid #42b983;
}

td {
  padding: 12px 15px;
  border-bottom: 1px solid #eee;
  font-size: 14px;
  vertical-align: middle;
}

.photo-cell { width: 70px; text-align: center; }

.student-img {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 50%;
  border: 2px solid #eee;
}

.no-photo {
  font-size: 10px;
  color: #999;
}

input, select {
  width: 100%;
  padding: 6px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.small-input { width: 50px; }

.center-text { text-align: center; }

.action-buttons {
  display: flex;
  gap: 8px;
}

button {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 600;
  transition: opacity 0.2s;
}

.edit-btn { background: #3498db; color: white; }
.delete-btn { background: #e74c3c; color: white; }
.save-btn { background: #42b983; color: white; }
.cancel-btn { background: #95a5a6; color: white; }

button:hover { opacity: 0.8; }

tr:hover { background-color: #fcfcfc; }
</style>