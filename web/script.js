const backupList = document.getElementById('backup-list');
const startBackupButton = document.getElementById('start-backup');

// Function to fetch backups
async function fetchBackups() {
  const response = await fetch('/api/backups');
  const data = await response.json();
  return data;
}

function displayBackups(backups) {
    backupList.innerHTML = '';
    backups.forEach(backup => {
      const backupElement = document.createElement('div');
      backupElement.classList.add('backup-card');
  
      // Create elements for each backup detail
      const backupId = document.createElement('p');
      const backupCreated = document.createElement('p');
      const backupType = document.createElement('p');
      const backupStatus = document.createElement('p');
      const statusSpan = document.createElement('span');
  
      // Set content for each element based on backup object properties
      backupId.textContent = `ID: ${backup.id}`;
      const createdDate = new Date(backup.created);
      backupCreated.textContent = `Created: ${createdDate.toLocaleDateString()}`;
      backupType.textContent = `Type: ${backup.type}`;
      backupStatus.textContent = 'Status: ';
      statusSpan.classList.add(`status-${backup.status.toLowerCase()}`); // Add class based on status
      statusSpan.textContent = backup.status;
      backupStatus.appendChild(statusSpan);
  
      // Append elements to the card
      backupElement.appendChild(backupId);
      backupElement.appendChild(backupCreated);
      backupElement.appendChild(backupType);
      backupElement.appendChild(backupStatus);
  
      // Add buttons with functionality (implementation details omitted for brevity)
      const deleteButton = document.createElement('button');
      deleteButton.textContent = 'Delete';
      deleteButton.addEventListener('click', () => {
        // Add logic to handle delete backup functionality using fetch
        console.log(`Delete Backup: ${backup.id}`); // Placeholder for deletion logic
      });
      const restoreButton = document.createElement('button');
      restoreButton.textContent = 'Restore';
      restoreButton.addEventListener('click', () => {
        // Add logic to handle restore backup functionality using fetch
        console.log(`Restore Backup: ${backup.id}`); // Placeholder for restoration logic
      });
      backupElement.appendChild(deleteButton);
      backupElement.appendChild(restoreButton);
  
      // Add the card to the backup list
      backupList.appendChild(backupElement);
    });
  }  

// Fetch backups on page load
fetchBackups().then(backups => displayBackups(backups));

// Add click event listener to start backup button
startBackupButton.addEventListener('click', async () => {

  // Add logic to handle start backup functionality using fetch
});
