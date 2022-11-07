using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class DataHandler : ICommandHandler<Data, DataResult>
    {
        private readonly IDataStore _dataStore;

        public DataHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<DataResult> ExecuteAsync(Data command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
