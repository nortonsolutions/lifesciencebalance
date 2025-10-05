/**
 * Admin Functions for course approval and management
 * @author Norton 2023
 */

// Load the admin dashboard with unapproved courses
export async function loadAdminDashboard() {
  // Get references to HTML elements
  const unapprovedCoursesContainer = document.getElementById('unapproved-courses');
  const approvedCoursesContainer = document.getElementById('approved-courses');
  const departmentFilter = document.getElementById('department-filter');
  const departmentFilterBtn = document.getElementById('filter-by-department-btn');

  // Load unapproved courses
  await loadUnapprovedCourses(unapprovedCoursesContainer);
  
  // Load approved courses
  await loadApprovedCourses(approvedCoursesContainer);

  // Setup event listeners for the department filter
  if (departmentFilterBtn) {
    departmentFilterBtn.addEventListener('click', async () => {
      const department = departmentFilter.value;
      if (department) {
        await loadCoursesByDepartment(approvedCoursesContainer, department);
      } else {
        // If no department specified, load all approved courses
        await loadApprovedCourses(approvedCoursesContainer);
      }
    });
  }
}

// Load unapproved courses
async function loadUnapprovedCourses(container) {
  if (!container) return;

  try {
    const response = await fetch('/course/unapproved');
    
    if (!response.ok) {
      throw new Error('Failed to fetch unapproved courses');
    }
    
    const courses = await response.json();
    
    // Clear the container
    container.innerHTML = '';
    
    if (courses.length === 0) {
      container.innerHTML = '<div class="alert alert-info">No courses waiting for approval.</div>';
      return;
    }
    
    // Create a card for each unapproved course
    courses.forEach(course => {
      const courseCard = createCourseCard(course, true);
      container.appendChild(courseCard);
    });
  } catch (error) {
    console.error('Error loading unapproved courses:', error);
    container.innerHTML = '<div class="alert alert-danger">Error loading unapproved courses. Please try again later.</div>';
  }
}

// Load approved courses
async function loadApprovedCourses(container) {
  if (!container) return;
  
  try {
    const response = await fetch('/course/approved');
    
    if (!response.ok) {
      throw new Error('Failed to fetch approved courses');
    }
    
    const courses = await response.json();
    
    // Clear the container
    container.innerHTML = '';
    
    if (courses.length === 0) {
      container.innerHTML = '<div class="alert alert-info">No approved courses found.</div>';
      return;
    }
    
    // Create a card for each approved course
    courses.forEach(course => {
      const courseCard = createCourseCard(course, false);
      container.appendChild(courseCard);
    });
  } catch (error) {
    console.error('Error loading approved courses:', error);
    container.innerHTML = '<div class="alert alert-danger">Error loading approved courses. Please try again later.</div>';
  }
}

// Load courses by department
async function loadCoursesByDepartment(container, department) {
  if (!container) return;
  
  try {
    const response = await fetch(`/course/department/${encodeURIComponent(department)}`);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch courses for department: ${department}`);
    }
    
    const courses = await response.json();
    
    // Clear the container
    container.innerHTML = '';
    
    if (courses.length === 0) {
      container.innerHTML = `<div class="alert alert-info">No courses found for department: ${department}</div>`;
      return;
    }
    
    // Create a card for each course in the department
    courses.forEach(course => {
      const courseCard = createCourseCard(course, false);
      container.appendChild(courseCard);
    });
  } catch (error) {
    console.error(`Error loading courses for department ${department}:`, error);
    container.innerHTML = '<div class="alert alert-danger">Error loading courses. Please try again later.</div>';
  }
}

// Create a course card element
function createCourseCard(course, isUnapproved) {
  const card = document.createElement('div');
  card.className = 'card mb-3';
  card.dataset.courseId = course.id;
  
  const cardBody = document.createElement('div');
  cardBody.className = 'card-body';
  
  // Course title
  const title = document.createElement('h5');
  title.className = 'card-title';
  title.textContent = course.name;
  
  // Department badge
  const department = document.createElement('span');
  department.className = 'badge bg-secondary me-2';
  department.textContent = course.department || 'No Department';
  
  // Description
  const description = document.createElement('p');
  description.className = 'card-text';
  description.textContent = course.description || 'No description available.';
  
  // Add elements to card body
  cardBody.appendChild(title);
  cardBody.appendChild(department);
  cardBody.appendChild(document.createElement('br'));
  cardBody.appendChild(description);
  
  // Add buttons for actions based on course status
  const buttonContainer = document.createElement('div');
  buttonContainer.className = 'd-flex gap-2';
  
  if (isUnapproved) {
    // Approve button for unapproved courses
    const approveBtn = document.createElement('button');
    approveBtn.className = 'btn btn-success';
    approveBtn.textContent = 'Approve Course';
    approveBtn.addEventListener('click', () => approveCourse(course.id));
    buttonContainer.appendChild(approveBtn);
  } else {
    // Unapprove button for approved courses
    const unapproveBtn = document.createElement('button');
    unapproveBtn.className = 'btn btn-warning';
    unapproveBtn.textContent = 'Unapprove Course';
    unapproveBtn.addEventListener('click', () => unapproveCourse(course.id));
    buttonContainer.appendChild(unapproveBtn);
  }
  
  // View details button for all courses
  const viewBtn = document.createElement('button');
  viewBtn.className = 'btn btn-primary';
  viewBtn.textContent = 'View Details';
  viewBtn.addEventListener('click', () => viewCourseDetails(course.id));
  buttonContainer.appendChild(viewBtn);
  
  cardBody.appendChild(buttonContainer);
  card.appendChild(cardBody);
  
  return card;
}

// Approve a course
async function approveCourse(courseId) {
  try {
    const response = await fetch(`/course/${courseId}/approve`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to approve course');
    }
    
    // Refresh the admin dashboard
    loadAdminDashboard();
    
    // Show success message
    showAlert('success', 'Course approved successfully.');
  } catch (error) {
    console.error('Error approving course:', error);
    showAlert('danger', 'Error approving course. Please try again.');
  }
}

// Unapprove a course
async function unapproveCourse(courseId) {
  try {
    const response = await fetch(`/course/${courseId}/unapprove`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to unapprove course');
    }
    
    // Refresh the admin dashboard
    loadAdminDashboard();
    
    // Show success message
    showAlert('success', 'Course unapproved.');
  } catch (error) {
    console.error('Error unapproving course:', error);
    showAlert('danger', 'Error unapproving course. Please try again.');
  }
}

// View course details
function viewCourseDetails(courseId) {
  window.location.href = `/course/${courseId}`;
}

// Show an alert message
function showAlert(type, message) {
  const alertContainer = document.getElementById('alert-container');
  if (!alertContainer) return;
  
  const alert = document.createElement('div');
  alert.className = `alert alert-${type} alert-dismissible fade show`;
  alert.role = 'alert';
  alert.innerHTML = `
    ${message}
    <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
  `;
  
  // Add the alert to the container
  alertContainer.appendChild(alert);
  
  // Automatically remove the alert after 5 seconds
  setTimeout(() => {
    alert.classList.remove('show');
    setTimeout(() => {
      alertContainer.removeChild(alert);
    }, 150);
  }, 5000);
}
