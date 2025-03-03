var yourNicknameKey = "yourNickname"

// Saving data to localStorage
function saveToLocalStorage(key, value) {
    try {
      // Convert the value to a JSON string before storing
      localStorage.setItem(key, JSON.stringify(value)); 
      console.log(`Successfully saved ${key} to localStorage.`);
    } catch (error) {
      console.error("Error saving to localStorage:", error);
      // Handle the error appropriately, e.g., display a message to the user
      // Common errors include exceeding storage quota or the user disabling storage
    }
  }
  

// Retrieving data from localStorage
function getFromLocalStorage(key) {
    try {
      const item = localStorage.getItem(key);
  
      // If the item exists, parse it back from a JSON string
      if (item) {
        return JSON.parse(item); 
      } else {
        return null; // Or undefined, depending on your preference
      }
    } catch (error) {
      console.error("Error retrieving from localStorage:", error);
      // Handle errors, e.g., invalid JSON
      return null; // Or handle differently
    }
  }