using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class HeloHandler : ICommandHandler<Helo, HeloResult>
    {
        private readonly IDataStore _dataStore;

        public HeloHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<HeloResult> ExecuteAsync(Helo command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
