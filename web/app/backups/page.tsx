'use client';

import { useEffect, useState } from "react";
import {Backup, ModelError as Error} from "../../components/spec";
import {NewBackupApi} from "../../components/api";
import {BackupCard, CreateBackup} from "./_components";

export default function Backups() {
  var [backups, setBackups] = useState([] as Backup[])
  var [error, setError] = useState<Error>()

  const loadBackups = () => {
    const api = NewBackupApi()
    api.listBackups()
      .then(data => {
        data.sort((a, b) => b.created.getDate() - a.created.getDate())
        setBackups(data)
      })
      .catch(err => {
        setError(err)
      })
  }

  useEffect(() => {
    loadBackups()
  }, [])

  return (
    <div>
      <p>{!error && error}</p>
      <CreateBackup refresh={loadBackups}/>
      <div className="flex flex-wrap gap-1 justify-center">
        {backups?.map(backup => <BackupCard 
          key={backup.id}
          backup={backup}
          onDelete={loadBackups}
        />)}
      </div>
    </div>
  )
}