/**
 * Course Submission Form functionality
 * @author Norton 2023
 */

// Initialize the course submission form
export function initCourseSubmissionForm() {
  const courseForm = document.getElementById('course-submission-form');
  
  if (!courseForm) return;
  
  courseForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    // Get form values
    const courseName = document.getElementById('course-name').value;
    const courseDescription = document.getElementById('course-description').value;
    const courseDepartment = document.getElementById('course-department').value;
    
    // Validate form
    if (!validateForm(courseName, courseDescription, courseDepartment)) {
      return;
    }
    
    // Get the current user ID from localStorage
    const userId = localStorage.getItem('userId');
    if (!userId) {
      showSubmissionMessage('error', 'You must be logged in to submit a course.');
      return;
    }
    
    // Create course object
    const newCourse = {
      name: courseName,
      description: courseDescription,
      department: courseDepartment,
      owner_id: parseInt(userId),
      approved: false,
      home_content: '',
      modules: []
    };
    
    try {
      // Submit course
      const response = await fetch('/course', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(newCourse)
      });
      
      if (!response.ok) {
        throw new Error('Failed to submit course');
      }
      
      const courseId = await response.json();
      
      // Create user-course relationship (making the user an instructor of the course)
      const userCourseResponse = await fetch('/usercourse', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          user_id: parseInt(userId),
          course_id: courseId,
          role_id: 2, // Assuming role_id 2 is "instructor"
          completion: 0,
          active: true
        })
      });
      
      if (!userCourseResponse.ok) {
        throw new Error('Failed to assign instructor role');
      }
      
      // Show success message
      showSubmissionMessage('success', 'Course submitted successfully! It will be reviewed by an administrator.');
      
      // Clear form
      courseForm.reset();
      
    } catch (error) {
      console.error('Error submitting course:', error);
      showSubmissionMessage('error', 'Error submitting course. Please try again later.');
    }
  });
}

// Validate form inputs
function validateForm(name, description, department) {
  if (!name || name.trim() === '') {
    showSubmissionMessage('error', 'Please enter a course name.');
    return false;
  }
  
  if (!description || description.trim() === '') {
    showSubmissionMessage('error', 'Please enter a course description.');
    return false;
  }
  
  if (!department || department.trim() === '') {
    showSubmissionMessage('error', 'Please select a department.');
    return false;
  }
  
  return true;
}

// Show a message after form submission
function showSubmissionMessage(type, message) {
  const messageContainer = document.getElementById('submission-message');
  if (!messageContainer) return;
  
  // Clear any existing messages
  messageContainer.innerHTML = '';
  
  // Create message element
  const messageElement = document.createElement('div');
  messageElement.className = type === 'success' 
    ? 'alert alert-success' 
    : 'alert alert-danger';
  messageElement.textContent = message;
  
  // Add message to the container
  messageContainer.appendChild(messageElement);
  
  // Scroll to the message
  messageElement.scrollIntoView({ behavior: 'smooth' });
  
  // Remove the message after 5 seconds if it's a success message
  if (type === 'success') {
    setTimeout(() => {
      messageContainer.innerHTML = '';
    }, 5000);
  }
}
