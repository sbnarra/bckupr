'use client';

import { useEffect, useState } from "react";
import {Backup, Error} from "../../components/spec";
import {NewBackupApi} from "../../components/api";
import {BackupCard} from "./_components";

export default function Backups() {
  var [backups, setBackups] = useState([] as Backup[])
  var [error, setError] = useState<Error>()

  useEffect(() => {
    NewBackupApi().listBackups(function(err: Error, data: [Backup], res: any) {
      console.log(res)
        if (error != null) {
            console.log("err:"+error)
            setError(err)
        } else {
          console.log("data:"+JSON.stringify(data))
          setBackups(data)
        }
    })
  }, [])

  return (
    <div>
      <p>{!error && error}</p>
        {backups?.map(backup => <BackupCard key={backup.id} backup={backup}/>)}
    </div>
  )
}