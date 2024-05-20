'use client';

import { useEffect, useState } from "react";
import {Button, Link} from '@nextui-org/react';
import {Backup, Error} from "../../components/spec";
import {NewBackupApi} from "../../components/api";
import {BackupCard, CreateBackup} from "./_components";

export default function Backups() {
  var [backups, setBackups] = useState([] as Backup[])
  var [error, setError] = useState<Error>()

  const loadBackups = () => {
    const api = NewBackupApi()
    api.listBackups(function(err: Error, data: [Backup], res: any) {
      console.log(res)
        if (error != null) {
          console.log("err:"+error)
          setError(err)
        } else {
          data.sort((a, b) => b.created - a.created)
          console.log("data:"+JSON.stringify(data))
          setBackups(data)
        }
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