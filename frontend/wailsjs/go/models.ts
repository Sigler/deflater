export namespace config {
	
	export class Alert {
	    package: string;
	    seen: string;
	
	    static createFrom(source: any = {}) {
	        return new Alert(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.package = source["package"];
	        this.seen = source["seen"];
	    }
	}
	export class Pending {
	    enable: string[];
	    disable: string[];
	    removeTasks?: string[];
	    token: string;
	    created: string;
	
	    static createFrom(source: any = {}) {
	        return new Pending(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable = source["enable"];
	        this.disable = source["disable"];
	        this.removeTasks = source["removeTasks"];
	        this.token = source["token"];
	        this.created = source["created"];
	    }
	}

}

export namespace main {
	
	export class FixResult {
	    id: string;
	    ok: boolean;
	    error?: string;
	    status: string;
	    phase: string;
	
	    static createFrom(source: any = {}) {
	        return new FixResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.ok = source["ok"];
	        this.error = source["error"];
	        this.status = source["status"];
	        this.phase = source["phase"];
	    }
	}
	export class ApplyOutcome {
	    needsElevation: boolean;
	    results: FixResult[];
	    saveWarning?: string;
	
	    static createFrom(source: any = {}) {
	        return new ApplyOutcome(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.needsElevation = source["needsElevation"];
	        this.results = this.convertValues(source["results"], FixResult);
	        this.saveWarning = source["saveWarning"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class RegOpInfo {
	    hive: string;
	    path: string;
	    name: string;
	    value: number;
	    revert: string;
	
	    static createFrom(source: any = {}) {
	        return new RegOpInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hive = source["hive"];
	        this.path = source["path"];
	        this.name = source["name"];
	        this.value = source["value"];
	        this.revert = source["revert"];
	    }
	}
	export class FixState {
	    id: string;
	    category: string;
	    kind: string;
	    caution: boolean;
	    profiles: string[];
	    reg?: RegOpInfo[];
	    appx?: string[];
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new FixState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.category = source["category"];
	        this.kind = source["kind"];
	        this.caution = source["caution"];
	        this.profiles = source["profiles"];
	        this.reg = this.convertValues(source["reg"], RegOpInfo);
	        this.appx = source["appx"];
	        this.status = source["status"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Report {
	    version: string;
	    elevated: boolean;
	    categories: string[];
	    fixes: FixState[];
	    managed: string[];
	    maintenance: boolean;
	    watcher: boolean;
	    alerts: config.Alert[];
	    taskMismatch: boolean;
	    conflictingTasks: schtask.ForeignTask[];
	    pending?: config.Pending;
	
	    static createFrom(source: any = {}) {
	        return new Report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.elevated = source["elevated"];
	        this.categories = source["categories"];
	        this.fixes = this.convertValues(source["fixes"], FixState);
	        this.managed = source["managed"];
	        this.maintenance = source["maintenance"];
	        this.watcher = source["watcher"];
	        this.alerts = this.convertValues(source["alerts"], config.Alert);
	        this.taskMismatch = source["taskMismatch"];
	        this.conflictingTasks = this.convertValues(source["conflictingTasks"], schtask.ForeignTask);
	        this.pending = this.convertValues(source["pending"], config.Pending);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ToggleResult {
	    saved: boolean;
	    needsElevation: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ToggleResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.saved = source["saved"];
	        this.needsElevation = source["needsElevation"];
	    }
	}

}

export namespace schtask {
	
	export class ForeignTask {
	    name: string;
	    tool: string;
	    note: string;
	
	    static createFrom(source: any = {}) {
	        return new ForeignTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.tool = source["tool"];
	        this.note = source["note"];
	    }
	}

}

export namespace update {
	
	export class Info {
	    available: boolean;
	    current: string;
	    latest: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.current = source["current"];
	        this.latest = source["latest"];
	        this.url = source["url"];
	    }
	}

}

