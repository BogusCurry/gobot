// specify the columns
var columnDefs = [
    {headerName:'UUID', field:'UUID', width:80, editable:true, sortable:true}, 
    {headerName:'Name', field:'Name', width:60, editable:true, sortable:true},
    {headerName:'Location', field:'Location', width:60, editable:true},
    {headerName:'Position', field:'Position', width:40, editable:true},
    {headerName:'Rotation', field:'Rotation', width:40, editable:true},
    {headerName:'Velocity', field:'Velocity', width:35, editable:true},
    {headerName:'LastUpdate', field:'LastUpdate', width:160, editable:true, sortable:true, sort:'desc'},
    {headerName:'OwnerName', field:'OwnerName', width:100, editable:true},
    {headerName:'OwnerKey', field:'OwnerKey', width:80, editable:true},
    {headerName:'ObjectType', field:'ObjectType', width:40, editable:true},
    {headerName:'ObjectClass', field:'ObjectClass', width:40, editable:true},
    {headerName:'RateEnergy', field:'RateEnergy', width:10, editable:true},
    {headerName:'RateMoney', field:'RateMoney', width:10, editable:true},      
    {headerName:'RateHappiness', field:'RateHappiness', width:10, editable:true},      
    {headerName:'PermURL', field:'PermURL', width:120, editable:false},
];

// get URLPathPrefix (hopefully it's already loaded by now)
var URLPathPrefix;

// let the grid know which columns and what data to use
var gridOptions = {
	columnDefs: columnDefs,
	rowData: null, // change to null since we're filling it from our own servers
	rowSelection: 'multiple',
	enableColResize: true,
	enableSorting: true,
	enableFilter: true,
	enableRangeSelection: true,
	suppressRowClickSelection: true,
	rowHeight: 22,
	animateRows: true,
	onModelUpdated: modelUpdated,
	debug: true,
	editType: 'fullRow',
    onRowValueChanged: function(event) {
        var data = event.data;
        console.log('onRowValueChanged: (' + data.UUID + ', ' + data.Name + ', ' + data.Location + ' ...)');
        
        var httpRequest = new XMLHttpRequest(); // see https://stackoverflow.com/questions/6418220/javascript-send-json-object-with-ajax
		httpRequest.open('POST', URLPathPrefix + '/uiPositionsUpdate/');
		httpRequest.setRequestHeader("Content-Type", "application/json");
		httpRequest.onreadystatechange = function() { // see https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest/onreadystatechange
			if(httpRequest.readyState === XMLHttpRequest.DONE && httpRequest.status === 200) {
				console.log(httpRequest.responseText);
  			}
  		}
  		var response = JSON.stringify(data)
  		// console.log('Response is going to be: ' + response)
  		httpRequest.send(response);
    },
	onGridReady: function() {
		gridOptions.api.sizeColumnsToFit()
	}
};


// wait for the document to be loaded, otherwise ag-Grid will not find the div in the document.
document.addEventListener("DOMContentLoaded", function() {
	URLPathPrefix = document.getElementById("URLPathPrefix").innerText;

	// lookup the container we want the Grid to use
	var eGridDiv = document.querySelector('#positionGrid');

	// create the grid passing in the div to use together with the columns & data we want to use
	new agGrid.Grid(eGridDiv, gridOptions);
	
	// do http request to get our sample data (see https://www.ag-grid.com/javascript-grid-value-getters/?framework=all#gsc.tab=0)
	var httpRequest = new XMLHttpRequest();
	httpRequest.open('GET', URLPathPrefix + '/uiPositions/');
	httpRequest.send();
	httpRequest.onreadystatechange = function() {
		if (httpRequest.readyState == 4 && httpRequest.status == 200) {
			var httpResult = JSON.parse(httpRequest.responseText);
			gridOptions.api.setRowData(httpResult);
		}
	};
});

function modelUpdated() {
}